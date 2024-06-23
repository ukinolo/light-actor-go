// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.0
// source: receiver.proto

package remote

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message  *anypb.Any `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Receiver string     `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reciver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_reciver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_reciver_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *Envelope) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reciver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_reciver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_reciver_proto_rawDescGZIP(), []int{1}
}

func (x *Error) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_reciver_proto protoreflect.FileDescriptor

var file_reciver_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x63, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x56, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x2e,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x22, 0x1d, 0x0a, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x43, 0x0a, 0x0e, 0x52, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x0e, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x2e,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x1a,
	0x0d, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x03,
	0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reciver_proto_rawDescOnce sync.Once
	file_reciver_proto_rawDescData = file_reciver_proto_rawDesc
)

func file_reciver_proto_rawDescGZIP() []byte {
	file_reciver_proto_rawDescOnce.Do(func() {
		file_reciver_proto_rawDescData = protoimpl.X.CompressGZIP(file_reciver_proto_rawDescData)
	})
	return file_reciver_proto_rawDescData
}

var file_reciver_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_reciver_proto_goTypes = []any{
	(*Envelope)(nil),  // 0: remote.Envelope
	(*Error)(nil),     // 1: remote.Error
	(*anypb.Any)(nil), // 2: google.protobuf.Any
}
var file_reciver_proto_depIdxs = []int32{
	2, // 0: remote.Envelope.message:type_name -> google.protobuf.Any
	0, // 1: remote.RemoteReceiver.ReceiveMessage:input_type -> remote.Envelope
	1, // 2: remote.RemoteReceiver.ReceiveMessage:output_type -> remote.Error
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_reciver_proto_init() }
func file_reciver_proto_init() {
	if File_reciver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_reciver_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Envelope); i {
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
		file_reciver_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_reciver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_reciver_proto_goTypes,
		DependencyIndexes: file_reciver_proto_depIdxs,
		MessageInfos:      file_reciver_proto_msgTypes,
	}.Build()
	File_reciver_proto = out.File
	file_reciver_proto_rawDesc = nil
	file_reciver_proto_goTypes = nil
	file_reciver_proto_depIdxs = nil
}
