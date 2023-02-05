package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/renbou/obzerva/internal/logs/storage/models"
)

var (
	jsonBoolStrs = map[bool]string{false: "false", true: "true"}
	jsonNullStr  = "null"
)

type flattenDecoder struct {
	*json.Decoder
	message models.FlatMessage
}

// flatten flattens a JSON message by unnesting all of its levels.
// Arrays elements are assigned their indices as their keys.
// Note that non-object/non-array messages are not mapped to any key,
// allowing only for non-scoped searches in the storage.
// If an error is encountered while iterating through the message,
// an empty FlatMessage is returned, just like for non-object/non-array messages.
func flatten(message []byte) models.FlatMessage {
	mr := bytes.NewReader(message)

	// UseNumber is needed to avoid screwing with the number representations
	decoder := flattenDecoder{Decoder: json.NewDecoder(mr)}
	decoder.UseNumber()

	if ok := decoder.flatten("", false); !ok {
		return models.FlatMessage{}
	}

	// Make sure there's no superfluous input left,
	// it's better to show an error somewhere later than to silently lose data
	_, err := decoder.Token()
	if err != io.EOF {
		return models.FlatMessage{}
	}

	// Sort fields by key to allow binary search for fast access later
	fields := decoder.message.Fields
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Key < fields[j].Key
	})

	return models.FlatMessage{Fields: fields}
}

// flatten here is called recursively after preparing the decoder,
// any error returned in this call or its descendants results in an empty FlatMessage
func (d *flattenDecoder) flatten(key string, nested bool) bool {
	token, err := d.Token()
	if err != nil {
		return false
	}

	switch t := token.(type) {
	case json.Delim:
		if t == '{' {
			return d.flattenObject(key)
		}

		// Delim here can never be } or ], because they are validated by decoder.Token first
		return d.flattenArray(key)
	case bool:
		d.appendIfNested(key, jsonBoolStrs[t], nested)
	case json.Number:
		d.appendIfNested(key, t.String(), nested)
	case string:
		d.appendIfNested(key, t, nested)
	case nil:
		d.appendIfNested(key, jsonNullStr, nested)
	}

	return true
}

func (d *flattenDecoder) appendIfNested(key, value string, nested bool) {
	// This additional check is needed to handle cases when a non-object/non-array value
	// is the root object, in which case nothing should be appended
	if !nested {
		return
	}

	d.message.Fields = append(d.message.Fields, &models.FlatMessage_KV{Key: key, Value: value})
}

func (d *flattenDecoder) flattenObject(keyPrefix string) bool {
	readObjectKey := func() (string, bool) {
		token, err := d.Token()
		if err != nil {
			return "", false
		}

		// json.Decoder guarantees that keys will be strings
		return token.(string), true
	}

	return d.flattenStructure(keyPrefix, readObjectKey, '}')
}

func (d *flattenDecoder) flattenArray(keyPrefix string) bool {
	var index int
	nextIndex := func() (string, bool) {
		cur := index
		index++

		// Itoa is already optimized, so for small arrays (logs shouldn't be massive, duh) this is fine
		return strconv.Itoa(cur), true
	}

	return d.flattenStructure(keyPrefix, nextIndex, ']')
}

func (d *flattenDecoder) flattenStructure(keyPrefix string, keyFunc func() (string, bool), endDelim rune) bool {
	for more := d.More(); more; more = d.More() {
		// Retrieve key. In case of an object, this will read a string token,
		// and in case of an array this will return the next index.
		key, ok := keyFunc()
		if !ok {
			return false
		}

		// Read structure value, which can be any nested object
		if ok := d.flatten(d.formatKey(keyPrefix, key), true); !ok {
			return false
		}
	}

	return d.assertDelim(endDelim)
}

func (d *flattenDecoder) assertDelim(delim rune) bool {
	token, err := d.Token()
	if err != nil || token != json.Delim(delim) {
		return false
	}

	return true
}

func (d *flattenDecoder) formatKey(prefix, key string) string {
	// replace invalid characters with _
	key = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' {
			return r
		}
		return '_'
	}, key)

	if prefix != "" {
		return prefix + "." + key
	}
	return key
}

// flatMessageToMapping returns a FlatMapping which searches through the FlatMessage using binary search.
func flatMessageToMapping(message *models.FlatMessage) FlatMapping {
	return func(key string) (value string, ok bool) {
		i, found := sort.Find(len(message.Fields), func(i int) int {
			return strings.Compare(key, message.Fields[i].Key)
		})

		if !found {
			return "", false
		}
		return message.Fields[i].Value, true
	}
}
