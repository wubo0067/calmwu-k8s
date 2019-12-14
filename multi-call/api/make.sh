#!/bin/bash

protoc -I=. -I=$GOPATH/src -I=../common/error --micro_out=. --go_out=. *.proto