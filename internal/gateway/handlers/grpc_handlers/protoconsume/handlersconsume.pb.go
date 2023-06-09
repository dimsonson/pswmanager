// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: handlersconsume.proto

package grpc_handlers

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

// PublishTextRecord
type ConsumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExchName      string `protobuf:"bytes,1,opt,name=exchName,proto3" json:"exchName,omitempty"`
	ConsumerQname string `protobuf:"bytes,2,opt,name=consumerQname,proto3" json:"consumerQname,omitempty"`
	RoutingKey    string `protobuf:"bytes,3,opt,name=routingKey,proto3" json:"routingKey,omitempty"`
}

func (x *ConsumeRequest) Reset() {
	*x = ConsumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handlersconsume_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRequest) ProtoMessage() {}

func (x *ConsumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_handlersconsume_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeRequest.ProtoReflect.Descriptor instead.
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return file_handlersconsume_proto_rawDescGZIP(), []int{0}
}

func (x *ConsumeRequest) GetExchName() string {
	if x != nil {
		return x.ExchName
	}
	return ""
}

func (x *ConsumeRequest) GetConsumerQname() string {
	if x != nil {
		return x.ConsumerQname
	}
	return ""
}

func (x *ConsumeRequest) GetRoutingKey() string {
	if x != nil {
		return x.RoutingKey
	}
	return ""
}

type ConsumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record     []byte `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
	RecordType int64  `protobuf:"varint,2,opt,name=recordType,proto3" json:"recordType,omitempty"`
	Error      string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"` // ошибка
}

func (x *ConsumeResponse) Reset() {
	*x = ConsumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handlersconsume_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeResponse) ProtoMessage() {}

func (x *ConsumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_handlersconsume_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeResponse.ProtoReflect.Descriptor instead.
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return file_handlersconsume_proto_rawDescGZIP(), []int{1}
}

func (x *ConsumeResponse) GetRecord() []byte {
	if x != nil {
		return x.Record
	}
	return nil
}

func (x *ConsumeResponse) GetRecordType() int64 {
	if x != nil {
		return x.RecordType
	}
	return 0
}

func (x *ConsumeResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_handlersconsume_proto protoreflect.FileDescriptor

var file_handlersconsume_proto_rawDesc = []byte{
	0x0a, 0x15, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x72,
	0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d,
	0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x51, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x51, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4b,
	0x65, 0x79, 0x22, 0x5f, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1e, 0x0a,
	0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x32, 0x4f, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x4d, 0x51,
	0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x12, 0x3a, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x30, 0x01, 0x42, 0x29, 0x5a, 0x27, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72,
	0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_handlersconsume_proto_rawDescOnce sync.Once
	file_handlersconsume_proto_rawDescData = file_handlersconsume_proto_rawDesc
)

func file_handlersconsume_proto_rawDescGZIP() []byte {
	file_handlersconsume_proto_rawDescOnce.Do(func() {
		file_handlersconsume_proto_rawDescData = protoimpl.X.CompressGZIP(file_handlersconsume_proto_rawDescData)
	})
	return file_handlersconsume_proto_rawDescData
}

var file_handlersconsume_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_handlersconsume_proto_goTypes = []interface{}{
	(*ConsumeRequest)(nil),  // 0: proto.ConsumeRequest
	(*ConsumeResponse)(nil), // 1: proto.ConsumeResponse
}
var file_handlersconsume_proto_depIdxs = []int32{
	0, // 0: proto.ServerRMQhandlers.Consume:input_type -> proto.ConsumeRequest
	1, // 1: proto.ServerRMQhandlers.Consume:output_type -> proto.ConsumeResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_handlersconsume_proto_init() }
func file_handlersconsume_proto_init() {
	if File_handlersconsume_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_handlersconsume_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeRequest); i {
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
		file_handlersconsume_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeResponse); i {
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
			RawDescriptor: file_handlersconsume_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_handlersconsume_proto_goTypes,
		DependencyIndexes: file_handlersconsume_proto_depIdxs,
		MessageInfos:      file_handlersconsume_proto_msgTypes,
	}.Build()
	File_handlersconsume_proto = out.File
	file_handlersconsume_proto_rawDesc = nil
	file_handlersconsume_proto_goTypes = nil
	file_handlersconsume_proto_depIdxs = nil
}
