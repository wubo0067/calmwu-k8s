FROM golang:alpine3.12 AS builder
WORKDIR /netlink

COPY go.mod .
COPY go.sum .
COPY testnetlink.go .
COPY createMacvlanEndpoint.go .

ENV GOPROXY https://goproxy.cn

RUN go mod vendor
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -o testnetlink testnetlink.go createMacvlanEndpoint.go