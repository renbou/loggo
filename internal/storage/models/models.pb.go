// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: internal/storage/proto/models.proto

// :TODO: use vtproto pool for fast unmarshaling

package models

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// FlatMessage is used for storing unnested log messages so that they can later be searched
// through in an optimal manner. Fields of the message are stored in a sorted order,
// allowing access in O(logN) without constructing a new map every time.
type FlatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields []*FlatMessage_KV `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *FlatMessage) Reset() {
	*x = FlatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_storage_proto_models_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlatMessage) ProtoMessage() {}

func (x *FlatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_storage_proto_models_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlatMessage.ProtoReflect.Descriptor instead.
func (*FlatMessage) Descriptor() ([]byte, []int) {
	return file_internal_storage_proto_models_proto_rawDescGZIP(), []int{0}
}

func (x *FlatMessage) GetFields() []*FlatMessage_KV {
	if x != nil {
		return x.Fields
	}
	return nil
}

// PreparedMessage is the message actually stored in the database, consisting of the original message,
// which will then be returned from the storage, and a flattened version which is used for fast search.
type PreparedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message []byte       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Flat    *FlatMessage `protobuf:"bytes,2,opt,name=flat,proto3" json:"flat,omitempty"`
}

func (x *PreparedMessage) Reset() {
	*x = PreparedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_storage_proto_models_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreparedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreparedMessage) ProtoMessage() {}

func (x *PreparedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_storage_proto_models_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreparedMessage.ProtoReflect.Descriptor instead.
func (*PreparedMessage) Descriptor() ([]byte, []int) {
	return file_internal_storage_proto_models_proto_rawDescGZIP(), []int{1}
}

func (x *PreparedMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *PreparedMessage) GetFlat() *FlatMessage {
	if x != nil {
		return x.Flat
	}
	return nil
}

type FlatMessage_KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FlatMessage_KV) Reset() {
	*x = FlatMessage_KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_storage_proto_models_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlatMessage_KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlatMessage_KV) ProtoMessage() {}

func (x *FlatMessage_KV) ProtoReflect() protoreflect.Message {
	mi := &file_internal_storage_proto_models_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlatMessage_KV.ProtoReflect.Descriptor instead.
func (*FlatMessage_KV) Descriptor() ([]byte, []int) {
	return file_internal_storage_proto_models_proto_rawDescGZIP(), []int{0, 0}
}

func (x *FlatMessage_KV) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FlatMessage_KV) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_internal_storage_proto_models_proto protoreflect.FileDescriptor

var file_internal_storage_proto_models_proto_rawDesc = []byte{
	0x0a, 0x23, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0x79, 0x0a, 0x0b, 0x46,
	0x6c, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3c, 0x0a, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c, 0x6f, 0x67,
	0x67, 0x6f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x46, 0x6c, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x56,
	0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x1a, 0x2c, 0x0a, 0x02, 0x4b, 0x56, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x62, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x46, 0x6c, 0x61, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x74, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x6e, 0x62, 0x6f, 0x75, 0x2f,
	0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_storage_proto_models_proto_rawDescOnce sync.Once
	file_internal_storage_proto_models_proto_rawDescData = file_internal_storage_proto_models_proto_rawDesc
)

func file_internal_storage_proto_models_proto_rawDescGZIP() []byte {
	file_internal_storage_proto_models_proto_rawDescOnce.Do(func() {
		file_internal_storage_proto_models_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_storage_proto_models_proto_rawDescData)
	})
	return file_internal_storage_proto_models_proto_rawDescData
}

var file_internal_storage_proto_models_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_storage_proto_models_proto_goTypes = []interface{}{
	(*FlatMessage)(nil),     // 0: loggo.storage.models.FlatMessage
	(*PreparedMessage)(nil), // 1: loggo.storage.models.PreparedMessage
	(*FlatMessage_KV)(nil),  // 2: loggo.storage.models.FlatMessage.KV
}
var file_internal_storage_proto_models_proto_depIdxs = []int32{
	2, // 0: loggo.storage.models.FlatMessage.fields:type_name -> loggo.storage.models.FlatMessage.KV
	0, // 1: loggo.storage.models.PreparedMessage.flat:type_name -> loggo.storage.models.FlatMessage
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_storage_proto_models_proto_init() }
func file_internal_storage_proto_models_proto_init() {
	if File_internal_storage_proto_models_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_storage_proto_models_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlatMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_storage_proto_models_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreparedMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_storage_proto_models_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlatMessage_KV); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_storage_proto_models_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_storage_proto_models_proto_goTypes,
		DependencyIndexes: file_internal_storage_proto_models_proto_depIdxs,
		MessageInfos:      file_internal_storage_proto_models_proto_msgTypes,
	}.Build()
	File_internal_storage_proto_models_proto = out.File
	file_internal_storage_proto_models_proto_rawDesc = nil
	file_internal_storage_proto_models_proto_goTypes = nil
	file_internal_storage_proto_models_proto_depIdxs = nil
}
