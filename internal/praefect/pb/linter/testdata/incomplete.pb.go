// Code generated by protoc-gen-go. DO NOT EDIT.
// source: incomplete.proto

package test

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// IncompleteRequest is missing the required option, so we expect a failure
type IncompleteRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IncompleteRequest) Reset()         { *m = IncompleteRequest{} }
func (m *IncompleteRequest) String() string { return proto.CompactTextString(m) }
func (*IncompleteRequest) ProtoMessage()    {}
func (*IncompleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_117c6bfd926ed1ba, []int{0}
}

func (m *IncompleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IncompleteRequest.Unmarshal(m, b)
}
func (m *IncompleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IncompleteRequest.Marshal(b, m, deterministic)
}
func (m *IncompleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IncompleteRequest.Merge(m, src)
}
func (m *IncompleteRequest) XXX_Size() int {
	return xxx_messageInfo_IncompleteRequest.Size(m)
}
func (m *IncompleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IncompleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IncompleteRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*IncompleteRequest)(nil), "test.IncompleteRequest")
}

func init() { proto.RegisterFile("incomplete.proto", fileDescriptor_117c6bfd926ed1ba) }

var fileDescriptor_117c6bfd926ed1ba = []byte{
	// 76 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0xcc, 0x4b, 0xce,
	0xcf, 0x2d, 0xc8, 0x49, 0x2d, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x49,
	0x2d, 0x2e, 0x91, 0xe2, 0xc9, 0x2f, 0x28, 0xa9, 0x2c, 0x80, 0x8a, 0x29, 0x09, 0x73, 0x09, 0x7a,
	0xc2, 0xd5, 0x05, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x24, 0xb1, 0x81, 0xe5, 0x8c, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xa0, 0x36, 0xd7, 0x43, 0x43, 0x00, 0x00, 0x00,
}
