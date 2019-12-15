// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/protobuf/common/sci_error/error.proto

package sci_v1_sci_error

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

type Error struct {
	Errno                int32    `protobuf:"varint,1,opt,name=errno,proto3" json:"errno,omitempty"`
	Domain               string   `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Reason               string   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Message              string   `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_64fb7fca21d6bea0, []int{0}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetErrno() int32 {
	if m != nil {
		return m.Errno
	}
	return 0
}

func (m *Error) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *Error) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Error)(nil), "sci.v1.sci_error.Error")
}

func init() {
	proto.RegisterFile("api/protobuf/common/sci_error/error.proto", fileDescriptor_64fb7fca21d6bea0)
}

var fileDescriptor_64fb7fca21d6bea0 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3,
	0x2f, 0x4e, 0xce, 0x8c, 0x4f, 0x2d, 0x2a, 0xca, 0x2f, 0xd2, 0x07, 0x93, 0x7a, 0x60, 0x79, 0x21,
	0x81, 0xe2, 0xe4, 0x4c, 0xbd, 0x32, 0x43, 0x3d, 0xb8, 0xac, 0x52, 0x3a, 0x17, 0xab, 0x2b, 0x88,
	0x21, 0x24, 0xc2, 0xc5, 0x9a, 0x5a, 0x54, 0x94, 0x97, 0x2f, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1a,
	0x04, 0xe1, 0x08, 0x89, 0x71, 0xb1, 0xa5, 0xe4, 0xe7, 0x26, 0x66, 0xe6, 0x49, 0x30, 0x29, 0x30,
	0x6a, 0x70, 0x06, 0x41, 0x79, 0x20, 0xf1, 0xa2, 0xd4, 0xc4, 0xe2, 0xfc, 0x3c, 0x09, 0x66, 0x88,
	0x38, 0x84, 0x27, 0x24, 0xc1, 0xc5, 0x9e, 0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x2a, 0xc1, 0x02,
	0x96, 0x80, 0x71, 0x93, 0xd8, 0xc0, 0x2e, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5e, 0xe6,
	0x36, 0xbb, 0xae, 0x00, 0x00, 0x00,
}