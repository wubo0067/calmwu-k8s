PROTO_SRCS	= $(wildcard *.proto)

.PHONY: all

all : proto

proto:
	protoc --proto_path=.:${GOPATH}/src --go_out=./ --micro_out=. $(PROTO_SRCS)