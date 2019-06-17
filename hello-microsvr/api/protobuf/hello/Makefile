PROTO_SRCS	= $(wildcard *.proto)

.PHONY: all

all : proto

proto:
	protoc --go_out=./ --micro_out=. $(PROTO_SRCS)