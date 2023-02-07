package storage

import (
	"fmt"

	"github.com/renbou/loggo/internal/storage/models"
	"google.golang.org/protobuf/proto"
)

func prepareMessage(m Message, flat *models.FlatMessage) ([]byte, error) {
	prepared := &models.PreparedMessage{
		Message: m,
		Flat:    flat,
	}

	data, err := proto.Marshal(prepared)
	if err != nil {
		return nil, fmt.Errorf("marshaling prepared message: %w", err)
	}
	return data, nil
}

func unprepareMessage(data []byte) (Message, *models.FlatMessage, error) {
	var prepared models.PreparedMessage

	if err := proto.Unmarshal(data, &prepared); err != nil {
		return nil, nil, fmt.Errorf("unmarshaling prepared message data: %w", err)
	}
	return prepared.Message, prepared.Flat, nil
}
