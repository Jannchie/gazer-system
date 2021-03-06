// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.2
// source: api/proto/server.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// Add a Task
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url        string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Tag        string `protobuf:"bytes,2,opt,name=tag,proto3" json:"tag,omitempty"`
	IntervalMS uint64 `protobuf:"varint,3,opt,name=intervalMS,proto3" json:"intervalMS,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{0}
}

func (x *Task) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Task) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Task) GetIntervalMS() uint64 {
	if x != nil {
		return x.IntervalMS
	}
	return 0
}

type Raw struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url       string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Tag       string                 `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Data      []byte                 `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Raw) Reset() {
	*x = Raw{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Raw) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Raw) ProtoMessage() {}

func (x *Raw) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Raw.ProtoReflect.Descriptor instead.
func (*Raw) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *Raw) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Raw) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Raw) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Raw) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Raw) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ListRawsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag   string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Url   string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Limit uint32 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListRawsReq) Reset() {
	*x = ListRawsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRawsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRawsReq) ProtoMessage() {}

func (x *ListRawsReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRawsReq.ProtoReflect.Descriptor instead.
func (*ListRawsReq) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{2}
}

func (x *ListRawsReq) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *ListRawsReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ListRawsReq) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type RawsResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Raws []*Raw `protobuf:"bytes,1,rep,name=raws,proto3" json:"raws,omitempty"`
}

func (x *RawsResp) Reset() {
	*x = RawsResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawsResp) ProtoMessage() {}

func (x *RawsResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawsResp.ProtoReflect.Descriptor instead.
func (*RawsResp) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{3}
}

func (x *RawsResp) GetRaws() []*Raw {
	if x != nil {
		return x.Raws
	}
	return nil
}

type ConsumeRawsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdList []uint64 `protobuf:"varint,1,rep,packed,name=idList,proto3" json:"idList,omitempty"`
}

func (x *ConsumeRawsReq) Reset() {
	*x = ConsumeRawsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRawsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRawsReq) ProtoMessage() {}

func (x *ConsumeRawsReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeRawsReq.ProtoReflect.Descriptor instead.
func (*ConsumeRawsReq) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{4}
}

func (x *ConsumeRawsReq) GetIdList() []uint64 {
	if x != nil {
		return x.IdList
	}
	return nil
}

type AddTasksReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks []*Task `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *AddTasksReq) Reset() {
	*x = AddTasksReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddTasksReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddTasksReq) ProtoMessage() {}

func (x *AddTasksReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddTasksReq.ProtoReflect.Descriptor instead.
func (*AddTasksReq) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{5}
}

func (x *AddTasksReq) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type OperationResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *OperationResp) Reset() {
	*x = OperationResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_server_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationResp) ProtoMessage() {}

func (x *OperationResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_server_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationResp.ProtoReflect.Descriptor instead.
func (*OperationResp) Descriptor() ([]byte, []int) {
	return file_api_proto_server_proto_rawDescGZIP(), []int{6}
}

func (x *OperationResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *OperationResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_api_proto_server_proto protoreflect.FileDescriptor

var file_api_proto_server_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x04, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x4d, 0x53, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x4d, 0x53, 0x22, 0x87, 0x01, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61,
	0x67, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x47, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x77, 0x73, 0x52, 0x65, 0x71, 0x12, 0x10,
	0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x24, 0x0a, 0x08, 0x52, 0x61, 0x77, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x04, 0x72, 0x61, 0x77, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x04, 0x2e, 0x52, 0x61, 0x77, 0x52, 0x04, 0x72, 0x61, 0x77, 0x73, 0x22, 0x28,
	0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x61, 0x77, 0x73, 0x52, 0x65, 0x71,
	0x12, 0x16, 0x0a, 0x06, 0x69, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04,
	0x52, 0x06, 0x69, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x2a, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x54,
	0x61, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74,
	0x61, 0x73, 0x6b, 0x73, 0x22, 0x35, 0x0a, 0x0d, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x8c, 0x01, 0x0a, 0x0b,
	0x47, 0x61, 0x7a, 0x65, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x28, 0x0a, 0x08, 0x41,
	0x64, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x0c, 0x2e, 0x41, 0x64, 0x64, 0x54, 0x61, 0x73,
	0x6b, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x23, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x77,
	0x73, 0x12, 0x0c, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x61, 0x77, 0x73, 0x52, 0x65, 0x71, 0x1a,
	0x09, 0x2e, 0x52, 0x61, 0x77, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2e, 0x0a, 0x0b, 0x43, 0x6f,
	0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x61, 0x77, 0x73, 0x12, 0x0f, 0x2e, 0x43, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x52, 0x61, 0x77, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x61,
	0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_server_proto_rawDescOnce sync.Once
	file_api_proto_server_proto_rawDescData = file_api_proto_server_proto_rawDesc
)

func file_api_proto_server_proto_rawDescGZIP() []byte {
	file_api_proto_server_proto_rawDescOnce.Do(func() {
		file_api_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_server_proto_rawDescData)
	})
	return file_api_proto_server_proto_rawDescData
}

var file_api_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_proto_server_proto_goTypes = []interface{}{
	(*Task)(nil),                  // 0: Task
	(*Raw)(nil),                   // 1: Raw
	(*ListRawsReq)(nil),           // 2: ListRawsReq
	(*RawsResp)(nil),              // 3: RawsResp
	(*ConsumeRawsReq)(nil),        // 4: ConsumeRawsReq
	(*AddTasksReq)(nil),           // 5: AddTasksReq
	(*OperationResp)(nil),         // 6: OperationResp
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_api_proto_server_proto_depIdxs = []int32{
	7, // 0: Raw.timestamp:type_name -> google.protobuf.Timestamp
	1, // 1: RawsResp.raws:type_name -> Raw
	0, // 2: AddTasksReq.tasks:type_name -> Task
	5, // 3: GazerSystem.AddTasks:input_type -> AddTasksReq
	2, // 4: GazerSystem.ListRaws:input_type -> ListRawsReq
	4, // 5: GazerSystem.ConsumeRaws:input_type -> ConsumeRawsReq
	6, // 6: GazerSystem.AddTasks:output_type -> OperationResp
	3, // 7: GazerSystem.ListRaws:output_type -> RawsResp
	6, // 8: GazerSystem.ConsumeRaws:output_type -> OperationResp
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_proto_server_proto_init() }
func file_api_proto_server_proto_init() {
	if File_api_proto_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_api_proto_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Raw); i {
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
		file_api_proto_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRawsReq); i {
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
		file_api_proto_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawsResp); i {
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
		file_api_proto_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeRawsReq); i {
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
		file_api_proto_server_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddTasksReq); i {
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
		file_api_proto_server_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperationResp); i {
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
			RawDescriptor: file_api_proto_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_server_proto_goTypes,
		DependencyIndexes: file_api_proto_server_proto_depIdxs,
		MessageInfos:      file_api_proto_server_proto_msgTypes,
	}.Build()
	File_api_proto_server_proto = out.File
	file_api_proto_server_proto_rawDesc = nil
	file_api_proto_server_proto_goTypes = nil
	file_api_proto_server_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GazerSystemClient is the client API for GazerSystem service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GazerSystemClient interface {
	AddTasks(ctx context.Context, in *AddTasksReq, opts ...grpc.CallOption) (*OperationResp, error)
	ListRaws(ctx context.Context, in *ListRawsReq, opts ...grpc.CallOption) (*RawsResp, error)
	ConsumeRaws(ctx context.Context, in *ConsumeRawsReq, opts ...grpc.CallOption) (*OperationResp, error)
}

type gazerSystemClient struct {
	cc grpc.ClientConnInterface
}

func NewGazerSystemClient(cc grpc.ClientConnInterface) GazerSystemClient {
	return &gazerSystemClient{cc}
}

func (c *gazerSystemClient) AddTasks(ctx context.Context, in *AddTasksReq, opts ...grpc.CallOption) (*OperationResp, error) {
	out := new(OperationResp)
	err := c.cc.Invoke(ctx, "/GazerSystem/AddTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gazerSystemClient) ListRaws(ctx context.Context, in *ListRawsReq, opts ...grpc.CallOption) (*RawsResp, error) {
	out := new(RawsResp)
	err := c.cc.Invoke(ctx, "/GazerSystem/ListRaws", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gazerSystemClient) ConsumeRaws(ctx context.Context, in *ConsumeRawsReq, opts ...grpc.CallOption) (*OperationResp, error) {
	out := new(OperationResp)
	err := c.cc.Invoke(ctx, "/GazerSystem/ConsumeRaws", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GazerSystemServer is the server API for GazerSystem service.
type GazerSystemServer interface {
	AddTasks(context.Context, *AddTasksReq) (*OperationResp, error)
	ListRaws(context.Context, *ListRawsReq) (*RawsResp, error)
	ConsumeRaws(context.Context, *ConsumeRawsReq) (*OperationResp, error)
}

// UnimplementedGazerSystemServer can be embedded to have forward compatible implementations.
type UnimplementedGazerSystemServer struct {
}

func (*UnimplementedGazerSystemServer) AddTasks(context.Context, *AddTasksReq) (*OperationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTasks not implemented")
}
func (*UnimplementedGazerSystemServer) ListRaws(context.Context, *ListRawsReq) (*RawsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRaws not implemented")
}
func (*UnimplementedGazerSystemServer) ConsumeRaws(context.Context, *ConsumeRawsReq) (*OperationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsumeRaws not implemented")
}

func RegisterGazerSystemServer(s *grpc.Server, srv GazerSystemServer) {
	s.RegisterService(&_GazerSystem_serviceDesc, srv)
}

func _GazerSystem_AddTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTasksReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GazerSystemServer).AddTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GazerSystem/AddTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GazerSystemServer).AddTasks(ctx, req.(*AddTasksReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GazerSystem_ListRaws_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRawsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GazerSystemServer).ListRaws(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GazerSystem/ListRaws",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GazerSystemServer).ListRaws(ctx, req.(*ListRawsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GazerSystem_ConsumeRaws_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsumeRawsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GazerSystemServer).ConsumeRaws(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GazerSystem/ConsumeRaws",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GazerSystemServer).ConsumeRaws(ctx, req.(*ConsumeRawsReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _GazerSystem_serviceDesc = grpc.ServiceDesc{
	ServiceName: "GazerSystem",
	HandlerType: (*GazerSystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTasks",
			Handler:    _GazerSystem_AddTasks_Handler,
		},
		{
			MethodName: "ListRaws",
			Handler:    _GazerSystem_ListRaws_Handler,
		},
		{
			MethodName: "ConsumeRaws",
			Handler:    _GazerSystem_ConsumeRaws_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/server.proto",
}
