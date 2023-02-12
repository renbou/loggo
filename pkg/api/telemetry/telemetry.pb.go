// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: api/telemetry.proto

package telemetry

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListLogMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From      *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To        *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Filter    *LogFilter             `protobuf:"bytes,3,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	PageSize  int32                  `protobuf:"varint,4,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string                 `protobuf:"bytes,5,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListLogMessagesRequest) Reset() {
	*x = ListLogMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLogMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLogMessagesRequest) ProtoMessage() {}

func (x *ListLogMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLogMessagesRequest.ProtoReflect.Descriptor instead.
func (*ListLogMessagesRequest) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{0}
}

func (x *ListLogMessagesRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *ListLogMessagesRequest) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *ListLogMessagesRequest) GetFilter() *LogFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *ListLogMessagesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListLogMessagesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListLogMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Batch *LogBatch `protobuf:"bytes,1,opt,name=batch,proto3" json:"batch,omitempty"`
}

func (x *ListLogMessagesResponse) Reset() {
	*x = ListLogMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLogMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLogMessagesResponse) ProtoMessage() {}

func (x *ListLogMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLogMessagesResponse.ProtoReflect.Descriptor instead.
func (*ListLogMessagesResponse) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{1}
}

func (x *ListLogMessagesResponse) GetBatch() *LogBatch {
	if x != nil {
		return x.Batch
	}
	return nil
}

type StreamLogMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From     *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	Filter   *LogFilter             `protobuf:"bytes,2,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	PageSize int32                  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *StreamLogMessagesRequest) Reset() {
	*x = StreamLogMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamLogMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamLogMessagesRequest) ProtoMessage() {}

func (x *StreamLogMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamLogMessagesRequest.ProtoReflect.Descriptor instead.
func (*StreamLogMessagesRequest) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{2}
}

func (x *StreamLogMessagesRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *StreamLogMessagesRequest) GetFilter() *LogFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *StreamLogMessagesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type StreamLogMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Response:
	//
	//	*StreamLogMessagesResponse_Batch
	//	*StreamLogMessagesResponse_Message
	Response isStreamLogMessagesResponse_Response `protobuf_oneof:"response"`
}

func (x *StreamLogMessagesResponse) Reset() {
	*x = StreamLogMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamLogMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamLogMessagesResponse) ProtoMessage() {}

func (x *StreamLogMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamLogMessagesResponse.ProtoReflect.Descriptor instead.
func (*StreamLogMessagesResponse) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{3}
}

func (m *StreamLogMessagesResponse) GetResponse() isStreamLogMessagesResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *StreamLogMessagesResponse) GetBatch() *LogBatch {
	if x, ok := x.GetResponse().(*StreamLogMessagesResponse_Batch); ok {
		return x.Batch
	}
	return nil
}

func (x *StreamLogMessagesResponse) GetMessage() []byte {
	if x, ok := x.GetResponse().(*StreamLogMessagesResponse_Message); ok {
		return x.Message
	}
	return nil
}

type isStreamLogMessagesResponse_Response interface {
	isStreamLogMessagesResponse_Response()
}

type StreamLogMessagesResponse_Batch struct {
	Batch *LogBatch `protobuf:"bytes,1,opt,name=batch,proto3,oneof"`
}

type StreamLogMessagesResponse_Message struct {
	Message []byte `protobuf:"bytes,2,opt,name=message,proto3,oneof"`
}

func (*StreamLogMessagesResponse_Batch) isStreamLogMessagesResponse_Response() {}

func (*StreamLogMessagesResponse_Message) isStreamLogMessagesResponse_Response() {}

type LogBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages      [][]byte `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	NextPageToken string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *LogBatch) Reset() {
	*x = LogBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogBatch) ProtoMessage() {}

func (x *LogBatch) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogBatch.ProtoReflect.Descriptor instead.
func (*LogBatch) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{4}
}

func (x *LogBatch) GetMessages() [][]byte {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *LogBatch) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type LogFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Filter:
	//
	//	*LogFilter_Text_
	//	*LogFilter_Scoped_
	//	*LogFilter_And_
	//	*LogFilter_Or_
	//	*LogFilter_Not_
	Filter isLogFilter_Filter `protobuf_oneof:"filter"`
}

func (x *LogFilter) Reset() {
	*x = LogFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter) ProtoMessage() {}

func (x *LogFilter) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter.ProtoReflect.Descriptor instead.
func (*LogFilter) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5}
}

func (m *LogFilter) GetFilter() isLogFilter_Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (x *LogFilter) GetText() *LogFilter_Text {
	if x, ok := x.GetFilter().(*LogFilter_Text_); ok {
		return x.Text
	}
	return nil
}

func (x *LogFilter) GetScoped() *LogFilter_Scoped {
	if x, ok := x.GetFilter().(*LogFilter_Scoped_); ok {
		return x.Scoped
	}
	return nil
}

func (x *LogFilter) GetAnd() *LogFilter_And {
	if x, ok := x.GetFilter().(*LogFilter_And_); ok {
		return x.And
	}
	return nil
}

func (x *LogFilter) GetOr() *LogFilter_Or {
	if x, ok := x.GetFilter().(*LogFilter_Or_); ok {
		return x.Or
	}
	return nil
}

func (x *LogFilter) GetNot() *LogFilter_Not {
	if x, ok := x.GetFilter().(*LogFilter_Not_); ok {
		return x.Not
	}
	return nil
}

type isLogFilter_Filter interface {
	isLogFilter_Filter()
}

type LogFilter_Text_ struct {
	Text *LogFilter_Text `protobuf:"bytes,1,opt,name=text,proto3,oneof"`
}

type LogFilter_Scoped_ struct {
	Scoped *LogFilter_Scoped `protobuf:"bytes,2,opt,name=scoped,proto3,oneof"`
}

type LogFilter_And_ struct {
	And *LogFilter_And `protobuf:"bytes,3,opt,name=and,proto3,oneof"`
}

type LogFilter_Or_ struct {
	Or *LogFilter_Or `protobuf:"bytes,4,opt,name=or,proto3,oneof"`
}

type LogFilter_Not_ struct {
	Not *LogFilter_Not `protobuf:"bytes,5,opt,name=not,proto3,oneof"`
}

func (*LogFilter_Text_) isLogFilter_Filter() {}

func (*LogFilter_Scoped_) isLogFilter_Filter() {}

func (*LogFilter_And_) isLogFilter_Filter() {}

func (*LogFilter_Or_) isLogFilter_Filter() {}

func (*LogFilter_Not_) isLogFilter_Filter() {}

type LogFilter_Text struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LogFilter_Text) Reset() {
	*x = LogFilter_Text{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter_Text) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter_Text) ProtoMessage() {}

func (x *LogFilter_Text) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter_Text.ProtoReflect.Descriptor instead.
func (*LogFilter_Text) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5, 0}
}

func (x *LogFilter_Text) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type LogFilter_Scoped struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LogFilter_Scoped) Reset() {
	*x = LogFilter_Scoped{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter_Scoped) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter_Scoped) ProtoMessage() {}

func (x *LogFilter_Scoped) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter_Scoped.ProtoReflect.Descriptor instead.
func (*LogFilter_Scoped) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5, 1}
}

func (x *LogFilter_Scoped) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *LogFilter_Scoped) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type LogFilter_And struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A *LogFilter `protobuf:"bytes,1,opt,name=a,proto3" json:"a,omitempty"`
	B *LogFilter `protobuf:"bytes,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *LogFilter_And) Reset() {
	*x = LogFilter_And{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter_And) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter_And) ProtoMessage() {}

func (x *LogFilter_And) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter_And.ProtoReflect.Descriptor instead.
func (*LogFilter_And) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5, 2}
}

func (x *LogFilter_And) GetA() *LogFilter {
	if x != nil {
		return x.A
	}
	return nil
}

func (x *LogFilter_And) GetB() *LogFilter {
	if x != nil {
		return x.B
	}
	return nil
}

type LogFilter_Or struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A *LogFilter `protobuf:"bytes,1,opt,name=a,proto3" json:"a,omitempty"`
	B *LogFilter `protobuf:"bytes,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *LogFilter_Or) Reset() {
	*x = LogFilter_Or{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter_Or) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter_Or) ProtoMessage() {}

func (x *LogFilter_Or) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter_Or.ProtoReflect.Descriptor instead.
func (*LogFilter_Or) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5, 3}
}

func (x *LogFilter_Or) GetA() *LogFilter {
	if x != nil {
		return x.A
	}
	return nil
}

func (x *LogFilter_Or) GetB() *LogFilter {
	if x != nil {
		return x.B
	}
	return nil
}

type LogFilter_Not struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A *LogFilter `protobuf:"bytes,1,opt,name=a,proto3" json:"a,omitempty"`
}

func (x *LogFilter_Not) Reset() {
	*x = LogFilter_Not{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_telemetry_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFilter_Not) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFilter_Not) ProtoMessage() {}

func (x *LogFilter_Not) ProtoReflect() protoreflect.Message {
	mi := &file_api_telemetry_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFilter_Not.ProtoReflect.Descriptor instead.
func (*LogFilter_Not) Descriptor() ([]byte, []int) {
	return file_api_telemetry_proto_rawDescGZIP(), []int{5, 4}
}

func (x *LogFilter_Not) GetA() *LogFilter {
	if x != nil {
		return x.A
	}
	return nil
}

var File_api_telemetry_proto protoreflect.FileDescriptor

var file_api_telemetry_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf8, 0x01, 0x0a, 0x16,
	0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x2a, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02,
	0x74, 0x6f, 0x12, 0x3b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x5f,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x4e, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x33, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c,
	0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x22, 0xaf, 0x01, 0x0a, 0x18, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x12, 0x3b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x7a, 0x0a, 0x19, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1a, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4e, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f,
	0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x84, 0x05, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x39, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x23, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c,
	0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x2e, 0x54, 0x65, 0x78, 0x74, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x3f, 0x0a,
	0x06, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e,
	0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65,
	0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x64, 0x48, 0x00, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x64, 0x12, 0x36,
	0x0a, 0x03, 0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c, 0x6f,
	0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x41, 0x6e, 0x64, 0x48,
	0x00, 0x52, 0x03, 0x61, 0x6e, 0x64, 0x12, 0x33, 0x0a, 0x02, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x2e, 0x4f, 0x72, 0x48, 0x00, 0x52, 0x02, 0x6f, 0x72, 0x12, 0x36, 0x0a, 0x03, 0x6e,
	0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c,
	0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x48, 0x00, 0x52, 0x03,
	0x6e, 0x6f, 0x74, 0x1a, 0x1c, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x1a, 0x34, 0x0a, 0x06, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x61, 0x0a, 0x03, 0x41, 0x6e, 0x64, 0x12, 0x2c,
	0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x67, 0x67,
	0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x01, 0x61, 0x12, 0x2c, 0x0a, 0x01,
	0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f,
	0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x01, 0x62, 0x1a, 0x60, 0x0a, 0x02, 0x4f, 0x72,
	0x12, 0x2c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f,
	0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x01, 0x61, 0x12, 0x2c,
	0x0a, 0x01, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x6f, 0x67, 0x67,
	0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x01, 0x62, 0x1a, 0x33, 0x0a, 0x03,
	0x4e, 0x6f, 0x74, 0x12, 0x2c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x01,
	0x61, 0x42, 0x08, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x32, 0xef, 0x01, 0x0a, 0x09,
	0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x12, 0x6c, 0x0a, 0x0f, 0x4c, 0x69, 0x73,
	0x74, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x2b, 0x2e, 0x6c,
	0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74,
	0x72, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x6c, 0x6f, 0x67, 0x67,
	0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x74, 0x0a, 0x11, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x2d, 0x2e, 0x6c,
	0x6f, 0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74,
	0x72, 0x79, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x6c, 0x6f,
	0x67, 0x67, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x2b, 0x5a,
	0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x6e, 0x62,
	0x6f, 0x75, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_api_telemetry_proto_rawDescOnce sync.Once
	file_api_telemetry_proto_rawDescData = file_api_telemetry_proto_rawDesc
)

func file_api_telemetry_proto_rawDescGZIP() []byte {
	file_api_telemetry_proto_rawDescOnce.Do(func() {
		file_api_telemetry_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_telemetry_proto_rawDescData)
	})
	return file_api_telemetry_proto_rawDescData
}

var file_api_telemetry_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_telemetry_proto_goTypes = []interface{}{
	(*ListLogMessagesRequest)(nil),    // 0: loggo.api.telemetry.ListLogMessagesRequest
	(*ListLogMessagesResponse)(nil),   // 1: loggo.api.telemetry.ListLogMessagesResponse
	(*StreamLogMessagesRequest)(nil),  // 2: loggo.api.telemetry.StreamLogMessagesRequest
	(*StreamLogMessagesResponse)(nil), // 3: loggo.api.telemetry.StreamLogMessagesResponse
	(*LogBatch)(nil),                  // 4: loggo.api.telemetry.LogBatch
	(*LogFilter)(nil),                 // 5: loggo.api.telemetry.LogFilter
	(*LogFilter_Text)(nil),            // 6: loggo.api.telemetry.LogFilter.Text
	(*LogFilter_Scoped)(nil),          // 7: loggo.api.telemetry.LogFilter.Scoped
	(*LogFilter_And)(nil),             // 8: loggo.api.telemetry.LogFilter.And
	(*LogFilter_Or)(nil),              // 9: loggo.api.telemetry.LogFilter.Or
	(*LogFilter_Not)(nil),             // 10: loggo.api.telemetry.LogFilter.Not
	(*timestamppb.Timestamp)(nil),     // 11: google.protobuf.Timestamp
}
var file_api_telemetry_proto_depIdxs = []int32{
	11, // 0: loggo.api.telemetry.ListLogMessagesRequest.from:type_name -> google.protobuf.Timestamp
	11, // 1: loggo.api.telemetry.ListLogMessagesRequest.to:type_name -> google.protobuf.Timestamp
	5,  // 2: loggo.api.telemetry.ListLogMessagesRequest.filter:type_name -> loggo.api.telemetry.LogFilter
	4,  // 3: loggo.api.telemetry.ListLogMessagesResponse.batch:type_name -> loggo.api.telemetry.LogBatch
	11, // 4: loggo.api.telemetry.StreamLogMessagesRequest.from:type_name -> google.protobuf.Timestamp
	5,  // 5: loggo.api.telemetry.StreamLogMessagesRequest.filter:type_name -> loggo.api.telemetry.LogFilter
	4,  // 6: loggo.api.telemetry.StreamLogMessagesResponse.batch:type_name -> loggo.api.telemetry.LogBatch
	6,  // 7: loggo.api.telemetry.LogFilter.text:type_name -> loggo.api.telemetry.LogFilter.Text
	7,  // 8: loggo.api.telemetry.LogFilter.scoped:type_name -> loggo.api.telemetry.LogFilter.Scoped
	8,  // 9: loggo.api.telemetry.LogFilter.and:type_name -> loggo.api.telemetry.LogFilter.And
	9,  // 10: loggo.api.telemetry.LogFilter.or:type_name -> loggo.api.telemetry.LogFilter.Or
	10, // 11: loggo.api.telemetry.LogFilter.not:type_name -> loggo.api.telemetry.LogFilter.Not
	5,  // 12: loggo.api.telemetry.LogFilter.And.a:type_name -> loggo.api.telemetry.LogFilter
	5,  // 13: loggo.api.telemetry.LogFilter.And.b:type_name -> loggo.api.telemetry.LogFilter
	5,  // 14: loggo.api.telemetry.LogFilter.Or.a:type_name -> loggo.api.telemetry.LogFilter
	5,  // 15: loggo.api.telemetry.LogFilter.Or.b:type_name -> loggo.api.telemetry.LogFilter
	5,  // 16: loggo.api.telemetry.LogFilter.Not.a:type_name -> loggo.api.telemetry.LogFilter
	0,  // 17: loggo.api.telemetry.Telemetry.ListLogMessages:input_type -> loggo.api.telemetry.ListLogMessagesRequest
	2,  // 18: loggo.api.telemetry.Telemetry.StreamLogMessages:input_type -> loggo.api.telemetry.StreamLogMessagesRequest
	1,  // 19: loggo.api.telemetry.Telemetry.ListLogMessages:output_type -> loggo.api.telemetry.ListLogMessagesResponse
	3,  // 20: loggo.api.telemetry.Telemetry.StreamLogMessages:output_type -> loggo.api.telemetry.StreamLogMessagesResponse
	19, // [19:21] is the sub-list for method output_type
	17, // [17:19] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_api_telemetry_proto_init() }
func file_api_telemetry_proto_init() {
	if File_api_telemetry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_telemetry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLogMessagesRequest); i {
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
		file_api_telemetry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLogMessagesResponse); i {
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
		file_api_telemetry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamLogMessagesRequest); i {
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
		file_api_telemetry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamLogMessagesResponse); i {
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
		file_api_telemetry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogBatch); i {
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
		file_api_telemetry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter); i {
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
		file_api_telemetry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter_Text); i {
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
		file_api_telemetry_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter_Scoped); i {
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
		file_api_telemetry_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter_And); i {
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
		file_api_telemetry_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter_Or); i {
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
		file_api_telemetry_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFilter_Not); i {
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
	file_api_telemetry_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_api_telemetry_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_api_telemetry_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*StreamLogMessagesResponse_Batch)(nil),
		(*StreamLogMessagesResponse_Message)(nil),
	}
	file_api_telemetry_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*LogFilter_Text_)(nil),
		(*LogFilter_Scoped_)(nil),
		(*LogFilter_And_)(nil),
		(*LogFilter_Or_)(nil),
		(*LogFilter_Not_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_telemetry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_telemetry_proto_goTypes,
		DependencyIndexes: file_api_telemetry_proto_depIdxs,
		MessageInfos:      file_api_telemetry_proto_msgTypes,
	}.Build()
	File_api_telemetry_proto = out.File
	file_api_telemetry_proto_rawDesc = nil
	file_api_telemetry_proto_goTypes = nil
	file_api_telemetry_proto_depIdxs = nil
}