// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: srv/stringprocess/stringprocess.proto

package eci_v1_svr_stringprocess

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for StringProcess service

type StringProcessService interface {
	ToUpper(ctx context.Context, in *OriginalStrReq, opts ...client.CallOption) (*UpperStrRes, error)
}

type stringProcessService struct {
	c    client.Client
	name string
}

func NewStringProcessService(name string, c client.Client) StringProcessService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "eci.v1.svr.stringprocess"
	}
	return &stringProcessService{
		c:    c,
		name: name,
	}
}

func (c *stringProcessService) ToUpper(ctx context.Context, in *OriginalStrReq, opts ...client.CallOption) (*UpperStrRes, error) {
	req := c.c.NewRequest(c.name, "StringProcess.ToUpper", in)
	out := new(UpperStrRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StringProcess service

type StringProcessHandler interface {
	ToUpper(context.Context, *OriginalStrReq, *UpperStrRes) error
}

func RegisterStringProcessHandler(s server.Server, hdlr StringProcessHandler, opts ...server.HandlerOption) error {
	type stringProcess interface {
		ToUpper(ctx context.Context, in *OriginalStrReq, out *UpperStrRes) error
	}
	type StringProcess struct {
		stringProcess
	}
	h := &stringProcessHandler{hdlr}
	return s.Handle(s.NewHandler(&StringProcess{h}, opts...))
}

type stringProcessHandler struct {
	StringProcessHandler
}

func (h *stringProcessHandler) ToUpper(ctx context.Context, in *OriginalStrReq, out *UpperStrRes) error {
	return h.StringProcessHandler.ToUpper(ctx, in, out)
}

// Client API for SplitProcess service

type SplitProcessService interface {
	Split(ctx context.Context, in *OriginalStrReq, opts ...client.CallOption) (*SplitStrRes, error)
}

type splitProcessService struct {
	c    client.Client
	name string
}

func NewSplitProcessService(name string, c client.Client) SplitProcessService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "eci.v1.svr.stringprocess"
	}
	return &splitProcessService{
		c:    c,
		name: name,
	}
}

func (c *splitProcessService) Split(ctx context.Context, in *OriginalStrReq, opts ...client.CallOption) (*SplitStrRes, error) {
	req := c.c.NewRequest(c.name, "SplitProcess.Split", in)
	out := new(SplitStrRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SplitProcess service

type SplitProcessHandler interface {
	Split(context.Context, *OriginalStrReq, *SplitStrRes) error
}

func RegisterSplitProcessHandler(s server.Server, hdlr SplitProcessHandler, opts ...server.HandlerOption) error {
	type splitProcess interface {
		Split(ctx context.Context, in *OriginalStrReq, out *SplitStrRes) error
	}
	type SplitProcess struct {
		splitProcess
	}
	h := &splitProcessHandler{hdlr}
	return s.Handle(s.NewHandler(&SplitProcess{h}, opts...))
}

type splitProcessHandler struct {
	SplitProcessHandler
}

func (h *splitProcessHandler) Split(ctx context.Context, in *OriginalStrReq, out *SplitStrRes) error {
	return h.SplitProcessHandler.Split(ctx, in, out)
}
