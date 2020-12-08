/*
 * @Author: CALM.WU
 * @Date: 2020-11-29 12:07:36
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-11-29 22:04:22
 */

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	protoHelloworld "istio-simplegrpc/proto/helloworld"

	"github.com/sanity-io/litter"
	calmwuUtils "github.com/wubo0067/calmwu-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

//
type GreeterServerImpl struct {
	// 这里必须嵌入，https://github.com/grpc/grpc-go/issues/3669
	protoHelloworld.UnimplementedGreeterServer
}

var (
	_index    = 0
	_hostName = ""
)

const (
	_defaultCallType = "CallType"
)

// SayHello 测试接口
func (gsi *GreeterServerImpl) SayHello(ctx context.Context, in *protoHelloworld.HelloRequest) (*protoHelloworld.HelloReply, error) {
	_index++

	// 获取http header
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(_defaultCallType)
	calmwuUtils.Debugf("index:%d Greeter.SayHello called, name: %s, CallType:%#v", _index, in.Name, vals)

	return &protoHelloworld.HelloReply{
		Message: fmt.Sprintf("srv-host:%s index:%d Hello %s", _hostName, _index, in.Name),
	}, nil
}

// CreateReservation 测试接口
func (gsi *GreeterServerImpl) CreateReservation(ctx context.Context, in *protoHelloworld.CreateReservationRequest) (*protoHelloworld.Reservation, error) {
	_index++

	// 获取http header
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(_defaultCallType)
	calmwuUtils.Debugf("index:%d Greeter.CreateReservation called, Reservation: %s, CallType:%#v", _index, litter.Sdump(in.Reservation), vals)

	return &protoHelloworld.Reservation{
		Id:        strconv.FormatInt(int64(_index), 10),
		Title:     fmt.Sprintf("%s-%d", in.Reservation.Title, _index),
		Venue:     in.Reservation.Venue,
		Room:      in.Reservation.Room,
		Timestamp: in.Reservation.Timestamp,
		Attendees: in.Reservation.Attendees,
	}, nil
}

var (
	_ protoHelloworld.GreeterServer = &GreeterServerImpl{}
)

func main() {
	calmwuUtils.Debug("istio-simplegrpc-server now start.")

	_hostName = os.Getenv("HOSTNAME")

	listen, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		calmwuUtils.Fatalf("failed to listen: %v", err.Error())
	}

	grpcSrv := grpc.NewServer()
	protoHelloworld.RegisterGreeterServer(grpcSrv, &GreeterServerImpl{})
	reflection.Register(grpcSrv)
	if err := grpcSrv.Serve(listen); err != nil {
		calmwuUtils.Fatal("failed to serve: %v", err.Error())
	}
}
