// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: app.proto

package app_pb

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

type RequestBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data        []byte  `protobuf:"bytes,1,opt,name=data,proto3,oneof" json:"data,omitempty"`
	WorkflowId  string  `protobuf:"bytes,2,opt,name=workflow_id,json=workflowId,proto3" json:"workflow_id,omitempty"`
	Depth       int32   `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`
	Width       int32   `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
	RequestType *string `protobuf:"bytes,5,opt,name=request_type,json=requestType,proto3,oneof" json:"request_type,omitempty"`
	PvPath      *string `protobuf:"bytes,6,opt,name=pv_path,json=pvPath,proto3,oneof" json:"pv_path,omitempty"`
}

func (x *RequestBody) Reset() {
	*x = RequestBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBody) ProtoMessage() {}

func (x *RequestBody) ProtoReflect() protoreflect.Message {
	mi := &file_app_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBody.ProtoReflect.Descriptor instead.
func (*RequestBody) Descriptor() ([]byte, []int) {
	return file_app_proto_rawDescGZIP(), []int{0}
}

func (x *RequestBody) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RequestBody) GetWorkflowId() string {
	if x != nil {
		return x.WorkflowId
	}
	return ""
}

func (x *RequestBody) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

func (x *RequestBody) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *RequestBody) GetRequestType() string {
	if x != nil && x.RequestType != nil {
		return *x.RequestType
	}
	return ""
}

func (x *RequestBody) GetPvPath() string {
	if x != nil && x.PvPath != nil {
		return *x.PvPath
	}
	return ""
}

type ResponseBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reply  *string `protobuf:"bytes,1,opt,name=reply,proto3,oneof" json:"reply,omitempty"`
	Code   int32   `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	PvPath *string `protobuf:"bytes,3,opt,name=pv_path,json=pvPath,proto3,oneof" json:"pv_path,omitempty"`
}

func (x *ResponseBody) Reset() {
	*x = ResponseBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseBody) ProtoMessage() {}

func (x *ResponseBody) ProtoReflect() protoreflect.Message {
	mi := &file_app_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseBody.ProtoReflect.Descriptor instead.
func (*ResponseBody) Descriptor() ([]byte, []int) {
	return file_app_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseBody) GetReply() string {
	if x != nil && x.Reply != nil {
		return *x.Reply
	}
	return ""
}

func (x *ResponseBody) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ResponseBody) GetPvPath() string {
	if x != nil && x.PvPath != nil {
		return *x.PvPath
	}
	return ""
}

var File_app_proto protoreflect.FileDescriptor

var file_app_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xdf, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x17, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x1f,
	0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x64, 0x65, 0x70, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x12, 0x26, 0x0a, 0x0c, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x70, 0x76, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x06, 0x70, 0x76, 0x50, 0x61, 0x74, 0x68, 0x88, 0x01,
	0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x70, 0x76, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x22, 0x71, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x19, 0x0a, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x88,
	0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x07, 0x70, 0x76, 0x5f, 0x70, 0x61, 0x74,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x06, 0x70, 0x76, 0x50, 0x61, 0x74,
	0x68, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x70, 0x76, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x32, 0x54, 0x0a, 0x0c, 0x67, 0x52,
	0x50, 0x43, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x44, 0x0a, 0x13, 0x67, 0x52,
	0x50, 0x43, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x72, 0x12, 0x15, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x1a, 0x16, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79,
	0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x61, 0x70, 0x70, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_proto_rawDescOnce sync.Once
	file_app_proto_rawDescData = file_app_proto_rawDesc
)

func file_app_proto_rawDescGZIP() []byte {
	file_app_proto_rawDescOnce.Do(func() {
		file_app_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_proto_rawDescData)
	})
	return file_app_proto_rawDescData
}

var file_app_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_proto_goTypes = []interface{}{
	(*RequestBody)(nil),  // 0: function.RequestBody
	(*ResponseBody)(nil), // 1: function.ResponseBody
}
var file_app_proto_depIdxs = []int32{
	0, // 0: function.gRPCFunction.gRPCFunctionHandler:input_type -> function.RequestBody
	1, // 1: function.gRPCFunction.gRPCFunctionHandler:output_type -> function.ResponseBody
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_app_proto_init() }
func file_app_proto_init() {
	if File_app_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestBody); i {
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
		file_app_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseBody); i {
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
	file_app_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_app_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_app_proto_goTypes,
		DependencyIndexes: file_app_proto_depIdxs,
		MessageInfos:      file_app_proto_msgTypes,
	}.Build()
	File_app_proto = out.File
	file_app_proto_rawDesc = nil
	file_app_proto_goTypes = nil
	file_app_proto_depIdxs = nil
}