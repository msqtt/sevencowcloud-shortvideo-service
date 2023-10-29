// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: user/activate_message.proto

package pb_usr

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

type ActivateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email        string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	ActivateCode string `protobuf:"bytes,2,opt,name=activate_code,json=activateCode,proto3" json:"activate_code,omitempty"`
}

func (x *ActivateUserRequest) Reset() {
	*x = ActivateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_activate_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateUserRequest) ProtoMessage() {}

func (x *ActivateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_activate_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateUserRequest.ProtoReflect.Descriptor instead.
func (*ActivateUserRequest) Descriptor() ([]byte, []int) {
	return file_user_activate_message_proto_rawDescGZIP(), []int{0}
}

func (x *ActivateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ActivateUserRequest) GetActivateCode() string {
	if x != nil {
		return x.ActivateCode
	}
	return ""
}

type ActivateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *ActivateUserResponse) Reset() {
	*x = ActivateUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_activate_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateUserResponse) ProtoMessage() {}

func (x *ActivateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_activate_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateUserResponse.ProtoReflect.Descriptor instead.
func (*ActivateUserResponse) Descriptor() ([]byte, []int) {
	return file_user_activate_message_proto_rawDescGZIP(), []int{1}
}

func (x *ActivateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_user_activate_message_proto protoreflect.FileDescriptor

var file_user_activate_message_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x13, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x39, 0x0a, 0x14, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x21, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x73, 0x71, 0x74, 0x74, 0x2f, 0x73, 0x65, 0x76, 0x65, 0x6e, 0x63, 0x6f,
	0x77, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2d, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62,
	0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x3b, 0x70, 0x62, 0x5f, 0x75, 0x73, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_activate_message_proto_rawDescOnce sync.Once
	file_user_activate_message_proto_rawDescData = file_user_activate_message_proto_rawDesc
)

func file_user_activate_message_proto_rawDescGZIP() []byte {
	file_user_activate_message_proto_rawDescOnce.Do(func() {
		file_user_activate_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_activate_message_proto_rawDescData)
	})
	return file_user_activate_message_proto_rawDescData
}

var file_user_activate_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_user_activate_message_proto_goTypes = []interface{}{
	(*ActivateUserRequest)(nil),  // 0: user.v1.ActivateUserRequest
	(*ActivateUserResponse)(nil), // 1: user.v1.ActivateUserResponse
	(*User)(nil),                 // 2: user.v1.User
}
var file_user_activate_message_proto_depIdxs = []int32{
	2, // 0: user.v1.ActivateUserResponse.user:type_name -> user.v1.User
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_user_activate_message_proto_init() }
func file_user_activate_message_proto_init() {
	if File_user_activate_message_proto != nil {
		return
	}
	file_user_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_user_activate_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateUserRequest); i {
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
		file_user_activate_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateUserResponse); i {
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
			RawDescriptor: file_user_activate_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_activate_message_proto_goTypes,
		DependencyIndexes: file_user_activate_message_proto_depIdxs,
		MessageInfos:      file_user_activate_message_proto_msgTypes,
	}.Build()
	File_user_activate_message_proto = out.File
	file_user_activate_message_proto_rawDesc = nil
	file_user_activate_message_proto_goTypes = nil
	file_user_activate_message_proto_depIdxs = nil
}
