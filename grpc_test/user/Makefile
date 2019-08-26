PROTO_SRCS	= $(wildcard *.proto)

.PHONY: all

all : proto

proto:
	protoc --go_out=plugins=grpc:./ $(PROTO_SRCS)