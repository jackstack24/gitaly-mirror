// Code generated by protoc-gen-go. DO NOT EDIT.
// source: go/internal/linter/testdata/invalid.proto

package test

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	gitalypb "gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type InvalidMethodRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InvalidMethodRequest) Reset()         { *m = InvalidMethodRequest{} }
func (m *InvalidMethodRequest) String() string { return proto.CompactTextString(m) }
func (*InvalidMethodRequest) ProtoMessage()    {}
func (*InvalidMethodRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{0}
}

func (m *InvalidMethodRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvalidMethodRequest.Unmarshal(m, b)
}
func (m *InvalidMethodRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvalidMethodRequest.Marshal(b, m, deterministic)
}
func (m *InvalidMethodRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvalidMethodRequest.Merge(m, src)
}
func (m *InvalidMethodRequest) XXX_Size() int {
	return xxx_messageInfo_InvalidMethodRequest.Size(m)
}
func (m *InvalidMethodRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InvalidMethodRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InvalidMethodRequest proto.InternalMessageInfo

type InvalidMethodRequestWithRepo struct {
	Destination          *gitalypb.Repository `protobuf:"bytes,1,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InvalidMethodRequestWithRepo) Reset()         { *m = InvalidMethodRequestWithRepo{} }
func (m *InvalidMethodRequestWithRepo) String() string { return proto.CompactTextString(m) }
func (*InvalidMethodRequestWithRepo) ProtoMessage()    {}
func (*InvalidMethodRequestWithRepo) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{1}
}

func (m *InvalidMethodRequestWithRepo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvalidMethodRequestWithRepo.Unmarshal(m, b)
}
func (m *InvalidMethodRequestWithRepo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvalidMethodRequestWithRepo.Marshal(b, m, deterministic)
}
func (m *InvalidMethodRequestWithRepo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvalidMethodRequestWithRepo.Merge(m, src)
}
func (m *InvalidMethodRequestWithRepo) XXX_Size() int {
	return xxx_messageInfo_InvalidMethodRequestWithRepo.Size(m)
}
func (m *InvalidMethodRequestWithRepo) XXX_DiscardUnknown() {
	xxx_messageInfo_InvalidMethodRequestWithRepo.DiscardUnknown(m)
}

var xxx_messageInfo_InvalidMethodRequestWithRepo proto.InternalMessageInfo

func (m *InvalidMethodRequestWithRepo) GetDestination() *gitalypb.Repository {
	if m != nil {
		return m.Destination
	}
	return nil
}

type InvalidTargetType struct {
	WrongType            int32    `protobuf:"varint,1,opt,name=wrong_type,json=wrongType,proto3" json:"wrong_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InvalidTargetType) Reset()         { *m = InvalidTargetType{} }
func (m *InvalidTargetType) String() string { return proto.CompactTextString(m) }
func (*InvalidTargetType) ProtoMessage()    {}
func (*InvalidTargetType) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{2}
}

func (m *InvalidTargetType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvalidTargetType.Unmarshal(m, b)
}
func (m *InvalidTargetType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvalidTargetType.Marshal(b, m, deterministic)
}
func (m *InvalidTargetType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvalidTargetType.Merge(m, src)
}
func (m *InvalidTargetType) XXX_Size() int {
	return xxx_messageInfo_InvalidTargetType.Size(m)
}
func (m *InvalidTargetType) XXX_DiscardUnknown() {
	xxx_messageInfo_InvalidTargetType.DiscardUnknown(m)
}

var xxx_messageInfo_InvalidTargetType proto.InternalMessageInfo

func (m *InvalidTargetType) GetWrongType() int32 {
	if m != nil {
		return m.WrongType
	}
	return 0
}

type InvalidMethodResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InvalidMethodResponse) Reset()         { *m = InvalidMethodResponse{} }
func (m *InvalidMethodResponse) String() string { return proto.CompactTextString(m) }
func (*InvalidMethodResponse) ProtoMessage()    {}
func (*InvalidMethodResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{3}
}

func (m *InvalidMethodResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvalidMethodResponse.Unmarshal(m, b)
}
func (m *InvalidMethodResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvalidMethodResponse.Marshal(b, m, deterministic)
}
func (m *InvalidMethodResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvalidMethodResponse.Merge(m, src)
}
func (m *InvalidMethodResponse) XXX_Size() int {
	return xxx_messageInfo_InvalidMethodResponse.Size(m)
}
func (m *InvalidMethodResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InvalidMethodResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InvalidMethodResponse proto.InternalMessageInfo

type InvalidNestedRequest struct {
	InnerMessage         *InvalidTargetType `protobuf:"bytes,1,opt,name=inner_message,json=innerMessage,proto3" json:"inner_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *InvalidNestedRequest) Reset()         { *m = InvalidNestedRequest{} }
func (m *InvalidNestedRequest) String() string { return proto.CompactTextString(m) }
func (*InvalidNestedRequest) ProtoMessage()    {}
func (*InvalidNestedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{4}
}

func (m *InvalidNestedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InvalidNestedRequest.Unmarshal(m, b)
}
func (m *InvalidNestedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InvalidNestedRequest.Marshal(b, m, deterministic)
}
func (m *InvalidNestedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InvalidNestedRequest.Merge(m, src)
}
func (m *InvalidNestedRequest) XXX_Size() int {
	return xxx_messageInfo_InvalidNestedRequest.Size(m)
}
func (m *InvalidNestedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InvalidNestedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InvalidNestedRequest proto.InternalMessageInfo

func (m *InvalidNestedRequest) GetInnerMessage() *InvalidTargetType {
	if m != nil {
		return m.InnerMessage
	}
	return nil
}

type RequestWithStorage struct {
	StorageName          string               `protobuf:"bytes,1,opt,name=storage_name,json=storageName,proto3" json:"storage_name,omitempty"`
	Destination          *gitalypb.Repository `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RequestWithStorage) Reset()         { *m = RequestWithStorage{} }
func (m *RequestWithStorage) String() string { return proto.CompactTextString(m) }
func (*RequestWithStorage) ProtoMessage()    {}
func (*RequestWithStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{5}
}

func (m *RequestWithStorage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithStorage.Unmarshal(m, b)
}
func (m *RequestWithStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithStorage.Marshal(b, m, deterministic)
}
func (m *RequestWithStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithStorage.Merge(m, src)
}
func (m *RequestWithStorage) XXX_Size() int {
	return xxx_messageInfo_RequestWithStorage.Size(m)
}
func (m *RequestWithStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithStorage.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithStorage proto.InternalMessageInfo

func (m *RequestWithStorage) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

func (m *RequestWithStorage) GetDestination() *gitalypb.Repository {
	if m != nil {
		return m.Destination
	}
	return nil
}

type RequestWithStorageAndRepo struct {
	StorageName          string               `protobuf:"bytes,1,opt,name=storage_name,json=storageName,proto3" json:"storage_name,omitempty"`
	Destination          *gitalypb.Repository `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RequestWithStorageAndRepo) Reset()         { *m = RequestWithStorageAndRepo{} }
func (m *RequestWithStorageAndRepo) String() string { return proto.CompactTextString(m) }
func (*RequestWithStorageAndRepo) ProtoMessage()    {}
func (*RequestWithStorageAndRepo) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{6}
}

func (m *RequestWithStorageAndRepo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithStorageAndRepo.Unmarshal(m, b)
}
func (m *RequestWithStorageAndRepo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithStorageAndRepo.Marshal(b, m, deterministic)
}
func (m *RequestWithStorageAndRepo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithStorageAndRepo.Merge(m, src)
}
func (m *RequestWithStorageAndRepo) XXX_Size() int {
	return xxx_messageInfo_RequestWithStorageAndRepo.Size(m)
}
func (m *RequestWithStorageAndRepo) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithStorageAndRepo.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithStorageAndRepo proto.InternalMessageInfo

func (m *RequestWithStorageAndRepo) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

func (m *RequestWithStorageAndRepo) GetDestination() *gitalypb.Repository {
	if m != nil {
		return m.Destination
	}
	return nil
}

type RequestWithNestedStorageAndRepo struct {
	InnerMessage         *RequestWithStorageAndRepo `protobuf:"bytes,1,opt,name=inner_message,json=innerMessage,proto3" json:"inner_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *RequestWithNestedStorageAndRepo) Reset()         { *m = RequestWithNestedStorageAndRepo{} }
func (m *RequestWithNestedStorageAndRepo) String() string { return proto.CompactTextString(m) }
func (*RequestWithNestedStorageAndRepo) ProtoMessage()    {}
func (*RequestWithNestedStorageAndRepo) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{7}
}

func (m *RequestWithNestedStorageAndRepo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithNestedStorageAndRepo.Unmarshal(m, b)
}
func (m *RequestWithNestedStorageAndRepo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithNestedStorageAndRepo.Marshal(b, m, deterministic)
}
func (m *RequestWithNestedStorageAndRepo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithNestedStorageAndRepo.Merge(m, src)
}
func (m *RequestWithNestedStorageAndRepo) XXX_Size() int {
	return xxx_messageInfo_RequestWithNestedStorageAndRepo.Size(m)
}
func (m *RequestWithNestedStorageAndRepo) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithNestedStorageAndRepo.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithNestedStorageAndRepo proto.InternalMessageInfo

func (m *RequestWithNestedStorageAndRepo) GetInnerMessage() *RequestWithStorageAndRepo {
	if m != nil {
		return m.InnerMessage
	}
	return nil
}

type RequestWithMultipleNestedStorage struct {
	InnerMessage         *RequestWithStorage `protobuf:"bytes,1,opt,name=inner_message,json=innerMessage,proto3" json:"inner_message,omitempty"`
	StorageName          string              `protobuf:"bytes,2,opt,name=storage_name,json=storageName,proto3" json:"storage_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *RequestWithMultipleNestedStorage) Reset()         { *m = RequestWithMultipleNestedStorage{} }
func (m *RequestWithMultipleNestedStorage) String() string { return proto.CompactTextString(m) }
func (*RequestWithMultipleNestedStorage) ProtoMessage()    {}
func (*RequestWithMultipleNestedStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{8}
}

func (m *RequestWithMultipleNestedStorage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithMultipleNestedStorage.Unmarshal(m, b)
}
func (m *RequestWithMultipleNestedStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithMultipleNestedStorage.Marshal(b, m, deterministic)
}
func (m *RequestWithMultipleNestedStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithMultipleNestedStorage.Merge(m, src)
}
func (m *RequestWithMultipleNestedStorage) XXX_Size() int {
	return xxx_messageInfo_RequestWithMultipleNestedStorage.Size(m)
}
func (m *RequestWithMultipleNestedStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithMultipleNestedStorage.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithMultipleNestedStorage proto.InternalMessageInfo

func (m *RequestWithMultipleNestedStorage) GetInnerMessage() *RequestWithStorage {
	if m != nil {
		return m.InnerMessage
	}
	return nil
}

func (m *RequestWithMultipleNestedStorage) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

type RequestWithInnerNestedStorage struct {
	Header               *RequestWithInnerNestedStorage_Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *RequestWithInnerNestedStorage) Reset()         { *m = RequestWithInnerNestedStorage{} }
func (m *RequestWithInnerNestedStorage) String() string { return proto.CompactTextString(m) }
func (*RequestWithInnerNestedStorage) ProtoMessage()    {}
func (*RequestWithInnerNestedStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{9}
}

func (m *RequestWithInnerNestedStorage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithInnerNestedStorage.Unmarshal(m, b)
}
func (m *RequestWithInnerNestedStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithInnerNestedStorage.Marshal(b, m, deterministic)
}
func (m *RequestWithInnerNestedStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithInnerNestedStorage.Merge(m, src)
}
func (m *RequestWithInnerNestedStorage) XXX_Size() int {
	return xxx_messageInfo_RequestWithInnerNestedStorage.Size(m)
}
func (m *RequestWithInnerNestedStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithInnerNestedStorage.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithInnerNestedStorage proto.InternalMessageInfo

func (m *RequestWithInnerNestedStorage) GetHeader() *RequestWithInnerNestedStorage_Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type RequestWithInnerNestedStorage_Header struct {
	StorageName          string   `protobuf:"bytes,1,opt,name=storage_name,json=storageName,proto3" json:"storage_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestWithInnerNestedStorage_Header) Reset()         { *m = RequestWithInnerNestedStorage_Header{} }
func (m *RequestWithInnerNestedStorage_Header) String() string { return proto.CompactTextString(m) }
func (*RequestWithInnerNestedStorage_Header) ProtoMessage()    {}
func (*RequestWithInnerNestedStorage_Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{9, 0}
}

func (m *RequestWithInnerNestedStorage_Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithInnerNestedStorage_Header.Unmarshal(m, b)
}
func (m *RequestWithInnerNestedStorage_Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithInnerNestedStorage_Header.Marshal(b, m, deterministic)
}
func (m *RequestWithInnerNestedStorage_Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithInnerNestedStorage_Header.Merge(m, src)
}
func (m *RequestWithInnerNestedStorage_Header) XXX_Size() int {
	return xxx_messageInfo_RequestWithInnerNestedStorage_Header.Size(m)
}
func (m *RequestWithInnerNestedStorage_Header) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithInnerNestedStorage_Header.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithInnerNestedStorage_Header proto.InternalMessageInfo

func (m *RequestWithInnerNestedStorage_Header) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

type RequestWithWrongTypeRepository struct {
	Header               *RequestWithWrongTypeRepository_Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *RequestWithWrongTypeRepository) Reset()         { *m = RequestWithWrongTypeRepository{} }
func (m *RequestWithWrongTypeRepository) String() string { return proto.CompactTextString(m) }
func (*RequestWithWrongTypeRepository) ProtoMessage()    {}
func (*RequestWithWrongTypeRepository) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{10}
}

func (m *RequestWithWrongTypeRepository) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithWrongTypeRepository.Unmarshal(m, b)
}
func (m *RequestWithWrongTypeRepository) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithWrongTypeRepository.Marshal(b, m, deterministic)
}
func (m *RequestWithWrongTypeRepository) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithWrongTypeRepository.Merge(m, src)
}
func (m *RequestWithWrongTypeRepository) XXX_Size() int {
	return xxx_messageInfo_RequestWithWrongTypeRepository.Size(m)
}
func (m *RequestWithWrongTypeRepository) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithWrongTypeRepository.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithWrongTypeRepository proto.InternalMessageInfo

func (m *RequestWithWrongTypeRepository) GetHeader() *RequestWithWrongTypeRepository_Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type RequestWithWrongTypeRepository_Header struct {
	Repository           *InvalidMethodResponse `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *RequestWithWrongTypeRepository_Header) Reset()         { *m = RequestWithWrongTypeRepository_Header{} }
func (m *RequestWithWrongTypeRepository_Header) String() string { return proto.CompactTextString(m) }
func (*RequestWithWrongTypeRepository_Header) ProtoMessage()    {}
func (*RequestWithWrongTypeRepository_Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{10, 0}
}

func (m *RequestWithWrongTypeRepository_Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithWrongTypeRepository_Header.Unmarshal(m, b)
}
func (m *RequestWithWrongTypeRepository_Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithWrongTypeRepository_Header.Marshal(b, m, deterministic)
}
func (m *RequestWithWrongTypeRepository_Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithWrongTypeRepository_Header.Merge(m, src)
}
func (m *RequestWithWrongTypeRepository_Header) XXX_Size() int {
	return xxx_messageInfo_RequestWithWrongTypeRepository_Header.Size(m)
}
func (m *RequestWithWrongTypeRepository_Header) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithWrongTypeRepository_Header.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithWrongTypeRepository_Header proto.InternalMessageInfo

func (m *RequestWithWrongTypeRepository_Header) GetRepository() *InvalidMethodResponse {
	if m != nil {
		return m.Repository
	}
	return nil
}

type RequestWithNestedRepoNotFlagged struct {
	Header               *RequestWithNestedRepoNotFlagged_Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *RequestWithNestedRepoNotFlagged) Reset()         { *m = RequestWithNestedRepoNotFlagged{} }
func (m *RequestWithNestedRepoNotFlagged) String() string { return proto.CompactTextString(m) }
func (*RequestWithNestedRepoNotFlagged) ProtoMessage()    {}
func (*RequestWithNestedRepoNotFlagged) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{11}
}

func (m *RequestWithNestedRepoNotFlagged) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged.Unmarshal(m, b)
}
func (m *RequestWithNestedRepoNotFlagged) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged.Marshal(b, m, deterministic)
}
func (m *RequestWithNestedRepoNotFlagged) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithNestedRepoNotFlagged.Merge(m, src)
}
func (m *RequestWithNestedRepoNotFlagged) XXX_Size() int {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged.Size(m)
}
func (m *RequestWithNestedRepoNotFlagged) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithNestedRepoNotFlagged.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithNestedRepoNotFlagged proto.InternalMessageInfo

func (m *RequestWithNestedRepoNotFlagged) GetHeader() *RequestWithNestedRepoNotFlagged_Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type RequestWithNestedRepoNotFlagged_Header struct {
	Repository           *gitalypb.Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RequestWithNestedRepoNotFlagged_Header) Reset() {
	*m = RequestWithNestedRepoNotFlagged_Header{}
}
func (m *RequestWithNestedRepoNotFlagged_Header) String() string { return proto.CompactTextString(m) }
func (*RequestWithNestedRepoNotFlagged_Header) ProtoMessage()    {}
func (*RequestWithNestedRepoNotFlagged_Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_506a53e91b227711, []int{11, 0}
}

func (m *RequestWithNestedRepoNotFlagged_Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header.Unmarshal(m, b)
}
func (m *RequestWithNestedRepoNotFlagged_Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header.Marshal(b, m, deterministic)
}
func (m *RequestWithNestedRepoNotFlagged_Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header.Merge(m, src)
}
func (m *RequestWithNestedRepoNotFlagged_Header) XXX_Size() int {
	return xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header.Size(m)
}
func (m *RequestWithNestedRepoNotFlagged_Header) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header.DiscardUnknown(m)
}

var xxx_messageInfo_RequestWithNestedRepoNotFlagged_Header proto.InternalMessageInfo

func (m *RequestWithNestedRepoNotFlagged_Header) GetRepository() *gitalypb.Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func init() {
	proto.RegisterType((*InvalidMethodRequest)(nil), "test.InvalidMethodRequest")
	proto.RegisterType((*InvalidMethodRequestWithRepo)(nil), "test.InvalidMethodRequestWithRepo")
	proto.RegisterType((*InvalidTargetType)(nil), "test.InvalidTargetType")
	proto.RegisterType((*InvalidMethodResponse)(nil), "test.InvalidMethodResponse")
	proto.RegisterType((*InvalidNestedRequest)(nil), "test.InvalidNestedRequest")
	proto.RegisterType((*RequestWithStorage)(nil), "test.RequestWithStorage")
	proto.RegisterType((*RequestWithStorageAndRepo)(nil), "test.RequestWithStorageAndRepo")
	proto.RegisterType((*RequestWithNestedStorageAndRepo)(nil), "test.RequestWithNestedStorageAndRepo")
	proto.RegisterType((*RequestWithMultipleNestedStorage)(nil), "test.RequestWithMultipleNestedStorage")
	proto.RegisterType((*RequestWithInnerNestedStorage)(nil), "test.RequestWithInnerNestedStorage")
	proto.RegisterType((*RequestWithInnerNestedStorage_Header)(nil), "test.RequestWithInnerNestedStorage.Header")
	proto.RegisterType((*RequestWithWrongTypeRepository)(nil), "test.RequestWithWrongTypeRepository")
	proto.RegisterType((*RequestWithWrongTypeRepository_Header)(nil), "test.RequestWithWrongTypeRepository.Header")
	proto.RegisterType((*RequestWithNestedRepoNotFlagged)(nil), "test.RequestWithNestedRepoNotFlagged")
	proto.RegisterType((*RequestWithNestedRepoNotFlagged_Header)(nil), "test.RequestWithNestedRepoNotFlagged.Header")
}

func init() {
	proto.RegisterFile("go/internal/linter/testdata/invalid.proto", fileDescriptor_506a53e91b227711)
}

var fileDescriptor_506a53e91b227711 = []byte{
	// 688 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x96, 0x5d, 0x4f, 0x13, 0x4d,
	0x14, 0xc7, 0x9f, 0x69, 0xfa, 0x34, 0x30, 0xe0, 0xdb, 0x44, 0x05, 0xeb, 0x0b, 0x64, 0xf1, 0x05,
	0x95, 0x6c, 0xa1, 0xa0, 0x22, 0xc1, 0x0b, 0x1a, 0x63, 0x44, 0x42, 0x13, 0x0b, 0x49, 0x13, 0x5f,
	0x52, 0x47, 0xf7, 0x64, 0xbb, 0xc9, 0x76, 0x77, 0xdd, 0x19, 0x30, 0xbd, 0xf3, 0xd2, 0x78, 0xe5,
	0x95, 0xf8, 0x1d, 0xfc, 0x02, 0xc6, 0x0f, 0xe0, 0x87, 0xe2, 0xca, 0xec, 0x4b, 0xcb, 0xee, 0xcc,
	0x2c, 0x8e, 0xc5, 0xbb, 0xcd, 0xf4, 0x9c, 0xff, 0xf9, 0xcd, 0x7f, 0xce, 0xc9, 0x29, 0xbe, 0x6d,
	0xfb, 0x35, 0xc7, 0xe3, 0x10, 0x7a, 0xd4, 0xad, 0xb9, 0xf1, 0x57, 0x8d, 0x03, 0xe3, 0x16, 0xe5,
	0xb4, 0xe6, 0x78, 0xfb, 0xd4, 0x75, 0x2c, 0x33, 0x08, 0x7d, 0xee, 0x93, 0x72, 0x74, 0x5e, 0x9d,
	0x64, 0x5d, 0x1a, 0x42, 0x7a, 0x66, 0x5c, 0xc4, 0xe7, 0x37, 0x93, 0xa0, 0x6d, 0xe0, 0x5d, 0xdf,
	0x6a, 0xc1, 0xfb, 0x3d, 0x60, 0xdc, 0x78, 0x81, 0xaf, 0xa8, 0xce, 0xdb, 0x0e, 0xef, 0xb6, 0x20,
	0xf0, 0xc9, 0x1a, 0x9e, 0xb0, 0x80, 0x71, 0xc7, 0xa3, 0xdc, 0xf1, 0xbd, 0x69, 0x34, 0x8b, 0xe6,
	0x27, 0xea, 0xc4, 0xb4, 0x1d, 0x4e, 0xdd, 0xbe, 0x19, 0x85, 0x30, 0x87, 0xfb, 0x61, 0xbf, 0x51,
	0xfe, 0xf6, 0x6b, 0x01, 0xb5, 0xb2, 0xc1, 0xc6, 0x2a, 0x3e, 0x97, 0x6a, 0xef, 0xd2, 0xd0, 0x06,
	0xbe, 0xdb, 0x0f, 0x80, 0xcc, 0x61, 0xfc, 0x21, 0xf4, 0x3d, 0xbb, 0xc3, 0xfb, 0x01, 0xc4, 0x7a,
	0xff, 0xa7, 0xb9, 0xe3, 0xf1, 0x79, 0x14, 0x64, 0x4c, 0xe1, 0x0b, 0x02, 0x15, 0x0b, 0x7c, 0x8f,
	0x81, 0xb1, 0x3b, 0xbc, 0x46, 0x13, 0x18, 0x87, 0x01, 0x2e, 0x59, 0xc7, 0xa7, 0x1c, 0xcf, 0x83,
	0xb0, 0xd3, 0x03, 0xc6, 0xa8, 0x0d, 0x29, 0xe8, 0x94, 0x19, 0x59, 0x61, 0x4a, 0x14, 0xad, 0xc9,
	0x38, 0x7a, 0x3b, 0x09, 0x36, 0x18, 0x26, 0x99, 0x7b, 0xef, 0x70, 0x3f, 0xa4, 0x36, 0x90, 0x5b,
	0x78, 0x92, 0x25, 0x9f, 0x1d, 0x8f, 0xf6, 0x12, 0xc9, 0xf1, 0x46, 0xf9, 0x53, 0x7c, 0xcf, 0xf4,
	0x97, 0x26, 0xed, 0x01, 0x59, 0xc9, 0x7b, 0x54, 0x2a, 0xf2, 0x28, 0xef, 0xce, 0x47, 0x84, 0x2f,
	0xc9, 0x55, 0x37, 0x3c, 0x2b, 0xf6, 0x5d, 0xbb, 0xf8, 0x9a, 0x66, 0x71, 0xd5, 0x03, 0xd9, 0x78,
	0x26, 0x43, 0x90, 0x38, 0x2a, 0x70, 0x3c, 0x56, 0x1b, 0x3b, 0x93, 0x18, 0x5b, 0xc8, 0x2f, 0x18,
	0xfc, 0x19, 0xe1, 0xd9, 0x4c, 0xec, 0xf6, 0x9e, 0xcb, 0x9d, 0xc0, 0x85, 0x5c, 0x45, 0xf2, 0x48,
	0x5d, 0x6a, 0xba, 0xa8, 0x54, 0xbe, 0x86, 0xe4, 0x58, 0xa9, 0xc0, 0x31, 0xe3, 0x2b, 0xc2, 0x57,
	0x33, 0x6a, 0x9b, 0x91, 0x48, 0x9e, 0xa4, 0x81, 0x2b, 0x5d, 0xa0, 0x16, 0x84, 0x29, 0xc2, 0x1d,
	0x09, 0x41, 0x4e, 0x32, 0x9f, 0xc6, 0x19, 0xad, 0x34, 0xb3, 0xba, 0x84, 0x2b, 0xc9, 0x89, 0xf6,
	0x53, 0x1a, 0x3f, 0x10, 0xbe, 0x96, 0xa9, 0xd1, 0x1e, 0x8c, 0xc3, 0xd1, 0x23, 0x92, 0x4d, 0x81,
	0xec, 0xae, 0x44, 0xa6, 0xc8, 0x4a, 0xd1, 0xd2, 0x0e, 0x18, 0x00, 0x6e, 0x0d, 0x01, 0x37, 0x30,
	0x0e, 0x87, 0xc1, 0xa9, 0xf0, 0xe5, 0xdc, 0xe4, 0xe4, 0xa7, 0xb0, 0x51, 0xfe, 0x12, 0x09, 0x65,
	0x92, 0x8c, 0xef, 0x48, 0xd1, 0x4a, 0x11, 0x41, 0xd3, 0xe7, 0x4f, 0x5c, 0x6a, 0xdb, 0x60, 0x91,
	0x67, 0x02, 0xfb, 0x82, 0xc4, 0xae, 0x4a, 0x53, 0xc3, 0xaf, 0x0f, 0xe1, 0xeb, 0x0a, 0x78, 0xd5,
	0xec, 0x65, 0xa2, 0xea, 0x3f, 0x31, 0x3e, 0x9d, 0xde, 0x6c, 0x07, 0xc2, 0x7d, 0xe7, 0x1d, 0x90,
	0xad, 0xe1, 0x49, 0x72, 0xd7, 0x45, 0x52, 0x55, 0x3a, 0x10, 0xb3, 0x56, 0x8f, 0x73, 0xc7, 0xf8,
	0x8f, 0x3c, 0x17, 0xc4, 0x96, 0x46, 0x17, 0xab, 0x1c, 0x1e, 0xcc, 0x97, 0xc6, 0x64, 0xc9, 0xfa,
	0x49, 0x25, 0x4b, 0xe4, 0xa5, 0x20, 0xb9, 0x4c, 0x8c, 0x62, 0xc9, 0xc1, 0x42, 0x38, 0x5e, 0x7a,
	0xec, 0xf0, 0x60, 0xbe, 0x3c, 0x86, 0xce, 0x22, 0x89, 0x77, 0xe5, 0xa4, 0xbc, 0x48, 0xe2, 0xbd,
	0x47, 0xae, 0xeb, 0x74, 0xbf, 0x9e, 0xf8, 0x2b, 0x41, 0xfc, 0x3e, 0xb9, 0xa1, 0xd5, 0x9e, 0x7a,
	0xea, 0x4d, 0x41, 0xfd, 0x01, 0x29, 0xda, 0x4c, 0x7a, 0x7a, 0xa2, 0xbb, 0xab, 0x82, 0xbb, 0xb9,
	0xe5, 0x38, 0x9a, 0xbb, 0x0f, 0xff, 0x5d, 0x37, 0x94, 0x48, 0x1b, 0x9f, 0xc9, 0x0f, 0xc4, 0x22,
	0xf9, 0xd3, 0x06, 0xd1, 0xeb, 0xe1, 0xd7, 0xa2, 0xf0, 0x52, 0xe1, 0xbb, 0xfd, 0xbd, 0x3c, 0x92,
	0xe5, 0xeb, 0x64, 0x4e, 0x63, 0x17, 0xe8, 0x0f, 0x89, 0x20, 0xbf, 0x3c, 0x62, 0x5f, 0x1c, 0x39,
	0xfd, 0x46, 0x94, 0x5c, 0x21, 0x37, 0x25, 0x62, 0xe5, 0xfe, 0xd5, 0xac, 0xf0, 0xb6, 0x12, 0xff,
	0xa1, 0x5c, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x3b, 0x0f, 0xf4, 0x8a, 0x91, 0x0a, 0x00, 0x00,
}
