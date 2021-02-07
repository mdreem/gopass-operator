// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: gopass_repository/repository.proto

package gopass_repository

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Authentication struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	SecretRef string `protobuf:"bytes,3,opt,name=secretRef,proto3" json:"secretRef,omitempty"`
	SecretKey string `protobuf:"bytes,4,opt,name=secretKey,proto3" json:"secretKey,omitempty"`
}

func (x *Authentication) Reset() {
	*x = Authentication{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gopass_repository_repository_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authentication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authentication) ProtoMessage() {}

func (x *Authentication) ProtoReflect() protoreflect.Message {
	mi := &file_gopass_repository_repository_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Authentication.ProtoReflect.Descriptor instead.
func (*Authentication) Descriptor() ([]byte, []int) {
	return file_gopass_repository_repository_proto_rawDescGZIP(), []int{0}
}

func (x *Authentication) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Authentication) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Authentication) GetSecretRef() string {
	if x != nil {
		return x.SecretRef
	}
	return ""
}

func (x *Authentication) GetSecretKey() string {
	if x != nil {
		return x.SecretKey
	}
	return ""
}

type Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RepositoryURL  string          `protobuf:"bytes,1,opt,name=repositoryURL,proto3" json:"repositoryURL,omitempty"`
	Authentication *Authentication `protobuf:"bytes,2,opt,name=authentication,proto3" json:"authentication,omitempty"`
}

func (x *Repository) Reset() {
	*x = Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gopass_repository_repository_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repository) ProtoMessage() {}

func (x *Repository) ProtoReflect() protoreflect.Message {
	mi := &file_gopass_repository_repository_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repository.ProtoReflect.Descriptor instead.
func (*Repository) Descriptor() ([]byte, []int) {
	return file_gopass_repository_repository_proto_rawDescGZIP(), []int{1}
}

func (x *Repository) GetRepositoryURL() string {
	if x != nil {
		return x.RepositoryURL
	}
	return ""
}

func (x *Repository) GetAuthentication() *Authentication {
	if x != nil {
		return x.Authentication
	}
	return nil
}

type RepositoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Successful   bool   `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *RepositoryResponse) Reset() {
	*x = RepositoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gopass_repository_repository_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepositoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepositoryResponse) ProtoMessage() {}

func (x *RepositoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gopass_repository_repository_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepositoryResponse.ProtoReflect.Descriptor instead.
func (*RepositoryResponse) Descriptor() ([]byte, []int) {
	return file_gopass_repository_repository_proto_rawDescGZIP(), []int{2}
}

func (x *RepositoryResponse) GetSuccessful() bool {
	if x != nil {
		return x.Successful
	}
	return false
}

func (x *RepositoryResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type Secret struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Secret) Reset() {
	*x = Secret{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gopass_repository_repository_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Secret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Secret) ProtoMessage() {}

func (x *Secret) ProtoReflect() protoreflect.Message {
	mi := &file_gopass_repository_repository_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Secret.ProtoReflect.Descriptor instead.
func (*Secret) Descriptor() ([]byte, []int) {
	return file_gopass_repository_repository_proto_rawDescGZIP(), []int{3}
}

func (x *Secret) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Secret) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type SecretList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secrets []*Secret `protobuf:"bytes,1,rep,name=secrets,proto3" json:"secrets,omitempty"`
}

func (x *SecretList) Reset() {
	*x = SecretList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gopass_repository_repository_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretList) ProtoMessage() {}

func (x *SecretList) ProtoReflect() protoreflect.Message {
	mi := &file_gopass_repository_repository_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretList.ProtoReflect.Descriptor instead.
func (*SecretList) Descriptor() ([]byte, []int) {
	return file_gopass_repository_repository_proto_rawDescGZIP(), []int{4}
}

func (x *SecretList) GetSecrets() []*Secret {
	if x != nil {
		return x.Secrets
	}
	return nil
}

var File_gopass_repository_repository_proto protoreflect.FileDescriptor

var file_gopass_repository_repository_proto_rawDesc = []byte{
	0x0a, 0x22, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x86, 0x01, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65,
	0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52,
	0x65, 0x66, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79,
	0x22, 0x7d, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x24,
	0x0a, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x55, 0x52, 0x4c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x55, 0x52, 0x4c, 0x12, 0x49, 0x0a, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67,
	0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x58, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x66, 0x75, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x66, 0x75, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x38, 0x0a, 0x06, 0x53, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x22, 0x41, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x33, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x07, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x32, 0xa4, 0x02, 0x0a, 0x11, 0x52, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5e, 0x0a, 0x14,
	0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x1a, 0x25, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x10,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x1a,
	0x25, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x11, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x41, 0x6c, 0x6c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x1d, 0x2e,
	0x67, 0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x1d, 0x2e, 0x67,
	0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0x21, 0x5a,
	0x1f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x67,
	0x6f, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gopass_repository_repository_proto_rawDescOnce sync.Once
	file_gopass_repository_repository_proto_rawDescData = file_gopass_repository_repository_proto_rawDesc
)

func file_gopass_repository_repository_proto_rawDescGZIP() []byte {
	file_gopass_repository_repository_proto_rawDescOnce.Do(func() {
		file_gopass_repository_repository_proto_rawDescData = protoimpl.X.CompressGZIP(file_gopass_repository_repository_proto_rawDescData)
	})
	return file_gopass_repository_repository_proto_rawDescData
}

var file_gopass_repository_repository_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_gopass_repository_repository_proto_goTypes = []interface{}{
	(*Authentication)(nil),     // 0: gopass_repository.Authentication
	(*Repository)(nil),         // 1: gopass_repository.Repository
	(*RepositoryResponse)(nil), // 2: gopass_repository.RepositoryResponse
	(*Secret)(nil),             // 3: gopass_repository.Secret
	(*SecretList)(nil),         // 4: gopass_repository.SecretList
}
var file_gopass_repository_repository_proto_depIdxs = []int32{
	0, // 0: gopass_repository.Repository.authentication:type_name -> gopass_repository.Authentication
	3, // 1: gopass_repository.SecretList.secrets:type_name -> gopass_repository.Secret
	1, // 2: gopass_repository.RepositoryService.InitializeRepository:input_type -> gopass_repository.Repository
	1, // 3: gopass_repository.RepositoryService.UpdateRepository:input_type -> gopass_repository.Repository
	1, // 4: gopass_repository.RepositoryService.FetchAllPasswords:input_type -> gopass_repository.Repository
	2, // 5: gopass_repository.RepositoryService.InitializeRepository:output_type -> gopass_repository.RepositoryResponse
	2, // 6: gopass_repository.RepositoryService.UpdateRepository:output_type -> gopass_repository.RepositoryResponse
	4, // 7: gopass_repository.RepositoryService.FetchAllPasswords:output_type -> gopass_repository.SecretList
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_gopass_repository_repository_proto_init() }
func file_gopass_repository_repository_proto_init() {
	if File_gopass_repository_repository_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gopass_repository_repository_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Authentication); i {
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
		file_gopass_repository_repository_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repository); i {
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
		file_gopass_repository_repository_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepositoryResponse); i {
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
		file_gopass_repository_repository_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Secret); i {
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
		file_gopass_repository_repository_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecretList); i {
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
			RawDescriptor: file_gopass_repository_repository_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gopass_repository_repository_proto_goTypes,
		DependencyIndexes: file_gopass_repository_repository_proto_depIdxs,
		MessageInfos:      file_gopass_repository_repository_proto_msgTypes,
	}.Build()
	File_gopass_repository_repository_proto = out.File
	file_gopass_repository_repository_proto_rawDesc = nil
	file_gopass_repository_repository_proto_goTypes = nil
	file_gopass_repository_repository_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RepositoryServiceClient is the client API for RepositoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RepositoryServiceClient interface {
	InitializeRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*RepositoryResponse, error)
	UpdateRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*RepositoryResponse, error)
	FetchAllPasswords(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*SecretList, error)
}

type repositoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRepositoryServiceClient(cc grpc.ClientConnInterface) RepositoryServiceClient {
	return &repositoryServiceClient{cc}
}

func (c *repositoryServiceClient) InitializeRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*RepositoryResponse, error) {
	out := new(RepositoryResponse)
	err := c.cc.Invoke(ctx, "/gopass_repository.RepositoryService/InitializeRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) UpdateRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*RepositoryResponse, error) {
	out := new(RepositoryResponse)
	err := c.cc.Invoke(ctx, "/gopass_repository.RepositoryService/UpdateRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) FetchAllPasswords(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*SecretList, error) {
	out := new(SecretList)
	err := c.cc.Invoke(ctx, "/gopass_repository.RepositoryService/FetchAllPasswords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepositoryServiceServer is the server API for RepositoryService service.
type RepositoryServiceServer interface {
	InitializeRepository(context.Context, *Repository) (*RepositoryResponse, error)
	UpdateRepository(context.Context, *Repository) (*RepositoryResponse, error)
	FetchAllPasswords(context.Context, *Repository) (*SecretList, error)
}

// UnimplementedRepositoryServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRepositoryServiceServer struct {
}

func (*UnimplementedRepositoryServiceServer) InitializeRepository(context.Context, *Repository) (*RepositoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitializeRepository not implemented")
}
func (*UnimplementedRepositoryServiceServer) UpdateRepository(context.Context, *Repository) (*RepositoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRepository not implemented")
}
func (*UnimplementedRepositoryServiceServer) FetchAllPasswords(context.Context, *Repository) (*SecretList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchAllPasswords not implemented")
}

func RegisterRepositoryServiceServer(s *grpc.Server, srv RepositoryServiceServer) {
	s.RegisterService(&_RepositoryService_serviceDesc, srv)
}

func _RepositoryService_InitializeRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Repository)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).InitializeRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gopass_repository.RepositoryService/InitializeRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).InitializeRepository(ctx, req.(*Repository))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_UpdateRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Repository)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).UpdateRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gopass_repository.RepositoryService/UpdateRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).UpdateRepository(ctx, req.(*Repository))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_FetchAllPasswords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Repository)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).FetchAllPasswords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gopass_repository.RepositoryService/FetchAllPasswords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).FetchAllPasswords(ctx, req.(*Repository))
	}
	return interceptor(ctx, in, info, handler)
}

var _RepositoryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gopass_repository.RepositoryService",
	HandlerType: (*RepositoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitializeRepository",
			Handler:    _RepositoryService_InitializeRepository_Handler,
		},
		{
			MethodName: "UpdateRepository",
			Handler:    _RepositoryService_UpdateRepository_Handler,
		},
		{
			MethodName: "FetchAllPasswords",
			Handler:    _RepositoryService_FetchAllPasswords_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gopass_repository/repository.proto",
}
