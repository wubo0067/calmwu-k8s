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

	protoHelloworld "istio-simplegrpc/proto/helloworld"

	calmwuUtils "github.com/wubo0067/calmwu-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//
type GreeterServerImpl struct{
	// 这里必须嵌入，https://github.com/grpc/grpc-go/issues/3669
	protoHelloworld.UnimplementedGreeterServer
}

var (
	_index = 0
	_hostName = ""
)

// 
func (gsi *GreeterServerImpl) SayHello(ctx context.Context, in *protoHelloworld.HelloRequest) (*protoHelloworld.HelloReply, error) {
	_index++
	calmwuUtils.Debugf("index:%d Greeter.SayHello called, name: %s", _index, in.Name)
	return &protoHelloworld.HelloReply{
		Message: fmt.Sprintf("srv-host:%s index:%d Hello %s", _hostName, _index, in.Name),
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
