// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/protobuf/sci_auth/sci_openapiauth.proto

// sci_v1_srv_opanapiauth 开放平台鉴权
package sci_v1_srv_opanapiauth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	sci_error "multi-call/api/protobuf/common/sci_error"
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

// AuthTokenReq
type AuthTokenReq struct {
	ReqID                string   `protobuf:"bytes,1,opt,name=reqID,proto3" json:"reqID,omitempty"`
	UserID               string   `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	UserPwd              string   `protobuf:"bytes,3,opt,name=userPwd,proto3" json:"userPwd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthTokenReq) Reset()         { *m = AuthTokenReq{} }
func (m *AuthTokenReq) String() string { return proto.CompactTextString(m) }
func (*AuthTokenReq) ProtoMessage()    {}
func (*AuthTokenReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a20687fc8208642a, []int{0}
}

func (m *AuthTokenReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthTokenReq.Unmarshal(m, b)
}
func (m *AuthTokenReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthTokenReq.Marshal(b, m, deterministic)
}
func (m *AuthTokenReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthTokenReq.Merge(m, src)
}
func (m *AuthTokenReq) XXX_Size() int {
	return xxx_messageInfo_AuthTokenReq.Size(m)
}
func (m *AuthTokenReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthTokenReq.DiscardUnknown(m)
}

var xxx_messageInfo_AuthTokenReq proto.InternalMessageInfo

func (m *AuthTokenReq) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *AuthTokenReq) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *AuthTokenReq) GetUserPwd() string {
	if m != nil {
		return m.UserPwd
	}
	return ""
}

type AuthTokenRes struct {
	ReqID                string           `protobuf:"bytes,1,opt,name=reqID,proto3" json:"reqID,omitempty"`
	UserID               string           `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	Token                string           `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Error                *sci_error.Error `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *AuthTokenRes) Reset()         { *m = AuthTokenRes{} }
func (m *AuthTokenRes) String() string { return proto.CompactTextString(m) }
func (*AuthTokenRes) ProtoMessage()    {}
func (*AuthTokenRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_a20687fc8208642a, []int{1}
}

func (m *AuthTokenRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthTokenRes.Unmarshal(m, b)
}
func (m *AuthTokenRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthTokenRes.Marshal(b, m, deterministic)
}
func (m *AuthTokenRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthTokenRes.Merge(m, src)
}
func (m *AuthTokenRes) XXX_Size() int {
	return xxx_messageInfo_AuthTokenRes.Size(m)
}
func (m *AuthTokenRes) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthTokenRes.DiscardUnknown(m)
}

var xxx_messageInfo_AuthTokenRes proto.InternalMessageInfo

func (m *AuthTokenRes) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *AuthTokenRes) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *AuthTokenRes) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AuthTokenRes) GetError() *sci_error.Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type VerifyTokenReq struct {
	ReqID                string   `protobuf:"bytes,1,opt,name=reqID,proto3" json:"reqID,omitempty"`
	JtwToken             string   `protobuf:"bytes,2,opt,name=jtwToken,proto3" json:"jtwToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyTokenReq) Reset()         { *m = VerifyTokenReq{} }
func (m *VerifyTokenReq) String() string { return proto.CompactTextString(m) }
func (*VerifyTokenReq) ProtoMessage()    {}
func (*VerifyTokenReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a20687fc8208642a, []int{2}
}

func (m *VerifyTokenReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyTokenReq.Unmarshal(m, b)
}
func (m *VerifyTokenReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyTokenReq.Marshal(b, m, deterministic)
}
func (m *VerifyTokenReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyTokenReq.Merge(m, src)
}
func (m *VerifyTokenReq) XXX_Size() int {
	return xxx_messageInfo_VerifyTokenReq.Size(m)
}
func (m *VerifyTokenReq) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyTokenReq.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyTokenReq proto.InternalMessageInfo

func (m *VerifyTokenReq) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *VerifyTokenReq) GetJtwToken() string {
	if m != nil {
		return m.JtwToken
	}
	return ""
}

type VerifyTokenRes struct {
	ReqID                string           `protobuf:"bytes,1,opt,name=reqID,proto3" json:"reqID,omitempty"`
	Error                *sci_error.Error `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *VerifyTokenRes) Reset()         { *m = VerifyTokenRes{} }
func (m *VerifyTokenRes) String() string { return proto.CompactTextString(m) }
func (*VerifyTokenRes) ProtoMessage()    {}
func (*VerifyTokenRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_a20687fc8208642a, []int{3}
}

func (m *VerifyTokenRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyTokenRes.Unmarshal(m, b)
}
func (m *VerifyTokenRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyTokenRes.Marshal(b, m, deterministic)
}
func (m *VerifyTokenRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyTokenRes.Merge(m, src)
}
func (m *VerifyTokenRes) XXX_Size() int {
	return xxx_messageInfo_VerifyTokenRes.Size(m)
}
func (m *VerifyTokenRes) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyTokenRes.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyTokenRes proto.InternalMessageInfo

func (m *VerifyTokenRes) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *VerifyTokenRes) GetError() *sci_error.Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func init() {
	proto.RegisterType((*AuthTokenReq)(nil), "sci.v1.srv.opanapiauth.AuthTokenReq")
	proto.RegisterType((*AuthTokenRes)(nil), "sci.v1.srv.opanapiauth.AuthTokenRes")
	proto.RegisterType((*VerifyTokenReq)(nil), "sci.v1.srv.opanapiauth.VerifyTokenReq")
	proto.RegisterType((*VerifyTokenRes)(nil), "sci.v1.srv.opanapiauth.VerifyTokenRes")
}

func init() {
	proto.RegisterFile("api/protobuf/sci_auth/sci_openapiauth.proto", fileDescriptor_a20687fc8208642a)
}

var fileDescriptor_a20687fc8208642a = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x41, 0x4f, 0x32, 0x31,
	0x10, 0x65, 0xbf, 0x4f, 0x50, 0x07, 0xc3, 0xa1, 0x21, 0xb8, 0xd9, 0xc4, 0x84, 0x6c, 0x8c, 0x21,
	0x31, 0x74, 0x23, 0xfa, 0x07, 0x34, 0x78, 0xe0, 0x24, 0x41, 0xe5, 0x28, 0x29, 0x6b, 0x09, 0xd5,
	0xdd, 0xb6, 0xb4, 0x5d, 0x88, 0xde, 0xfc, 0xa7, 0xfe, 0x14, 0xd3, 0x2e, 0x6b, 0xc0, 0x48, 0x84,
	0x4b, 0x33, 0x6f, 0xe6, 0xcd, 0xdb, 0x79, 0x33, 0x0b, 0xe7, 0x44, 0xb2, 0x48, 0x2a, 0x61, 0xc4,
	0x38, 0x9b, 0x44, 0x3a, 0x66, 0x23, 0x92, 0x99, 0xa9, 0x0b, 0x84, 0xa4, 0x9c, 0x48, 0x66, 0x31,
	0x76, 0x0c, 0xd4, 0xd0, 0x31, 0xc3, 0xf3, 0x0b, 0xac, 0xd5, 0x1c, 0x0b, 0x49, 0x8a, 0x6a, 0x70,
	0x95, 0x66, 0x89, 0x61, 0xed, 0x98, 0x24, 0x49, 0xb4, 0xa6, 0x17, 0x8b, 0x34, 0x15, 0xdc, 0xa9,
	0x51, 0xa5, 0x84, 0x8a, 0xdc, 0x9b, 0xab, 0x85, 0x43, 0x38, 0xba, 0xce, 0xcc, 0xf4, 0x41, 0xbc,
	0x52, 0x3e, 0xa0, 0x33, 0x54, 0x87, 0xb2, 0xa2, 0xb3, 0x5e, 0xd7, 0xf7, 0x9a, 0x5e, 0xeb, 0x70,
	0x90, 0x03, 0xd4, 0x80, 0x4a, 0xa6, 0xa9, 0xea, 0x75, 0xfd, 0x7f, 0x2e, 0xbd, 0x44, 0xc8, 0x87,
	0x7d, 0x1b, 0xf5, 0x17, 0xcf, 0xfe, 0x7f, 0x57, 0x28, 0x60, 0xf8, 0xe1, 0xad, 0x09, 0xeb, 0x1d,
	0x85, 0xeb, 0x50, 0x36, 0xb6, 0x73, 0x29, 0x9b, 0x03, 0xd4, 0x86, 0xb2, 0x9b, 0xdd, 0xdf, 0x6b,
	0x7a, 0xad, 0x6a, 0xe7, 0x18, 0x17, 0xab, 0x28, 0xac, 0xe1, 0x5b, 0xfb, 0x0e, 0x72, 0x56, 0x78,
	0x03, 0xb5, 0x21, 0x55, 0x6c, 0xf2, 0xf6, 0x87, 0xbb, 0x00, 0x0e, 0x5e, 0xcc, 0xc2, 0x91, 0x96,
	0x63, 0x7c, 0xe3, 0xf0, 0xf1, 0x87, 0xc6, 0x26, 0x23, 0xbb, 0x8d, 0xd6, 0xf9, 0xf4, 0xe0, 0xe4,
	0x4e, 0x52, 0xde, 0x4f, 0x88, 0x99, 0x08, 0x95, 0xda, 0x55, 0x51, 0x6e, 0x58, 0x4c, 0x0c, 0x13,
	0xfc, 0x5e, 0xcd, 0xd1, 0x13, 0xd4, 0x6c, 0x52, 0x28, 0xf6, 0x4e, 0xdd, 0xb7, 0xd1, 0x29, 0xfe,
	0xfd, 0xf2, 0x78, 0xf5, 0x80, 0xc1, 0x36, 0x2c, 0x1d, 0x96, 0xd0, 0x08, 0xaa, 0x2b, 0xc6, 0xd0,
	0xd9, 0xa6, 0xb6, 0xf5, 0x0d, 0x06, 0xdb, 0xf1, 0x74, 0x58, 0x1a, 0x57, 0xdc, 0x0f, 0x76, 0xf9,
	0x15, 0x00, 0x00, 0xff, 0xff, 0x9c, 0x55, 0x43, 0xfa, 0xdd, 0x02, 0x00, 0x00,
}