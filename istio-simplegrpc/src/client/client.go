/*
 * @Author: CALM.WU
 * @Date: 2020-11-29 11:46:02
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-11-29 12:03:07
 */

package main

import (
	"context"
	protoHelloworld "istio-simplegrpc/proto/helloworld"
	"os"
	"time"

	calmwuUtils "github.com/wubo0067/calmwu-go/utils"
	"google.golang.org/grpc"
)

const (
	_helloWorldServiceAddr = "helloworld.istio-ns.svc.cluster.local:8081"
	_defaultName = "CalmWU"
	_defaultGRPCTimeout = 10 * time.Second
)

func main() {
	conn, err := grpc.Dial(_helloWorldServiceAddr, grpc.WithInsecure())
	if err != nil {
		calmwuUtils.Fatalf("grpc dial to %s failed. err:%s", err.Error())
	}
	defer conn.Close()

	client := protoHelloworld.NewGreeterClient(conn)
	name := _defaultName

	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), _defaultGRPCTimeout)
	defer cancel()

	resp, err := client.SayHello(ctx, &protoHelloworld.HelloRequest{
		Name: name,})
	if err != nil {
		calmwuUtils.Errorf("call greet.SayHello failed. err:%s", err.Error())
	} else {
		calmwuUtils.Debugf("call greet.SayHello resp:%#v", resp)
	}
}