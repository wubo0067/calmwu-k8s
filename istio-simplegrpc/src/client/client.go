/*
 * @Author: CALM.WU
 * @Date: 2020-11-29 11:46:02
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-11-29 22:04:41
 */

package main

import (
	"context"
	protoHelloworld "istio-simplegrpc/proto/helloworld"
	"os"
	"os/signal"
	"syscall"
	"time"

	calmwuUtils "github.com/wubo0067/calmwu-go/utils"
	"google.golang.org/grpc"
)

const (
	_greeterServiceAddr = "greeter.istio-ns.svc.cluster.local:8081"
	_defaultName = "CalmWU"
	_defaultGRPCTimeout = 10 * time.Second
)

var (
	_index = 0
)

func main() {
	pCtx, pCancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		killSignal := <-interrupt
		switch killSignal {
		case os.Interrupt:
			calmwuUtils.Debug("Got SIGINT...")
			pCancel()
		case syscall.SIGTERM:
			// 容器推出会收到该信号
			calmwuUtils.Debug("Got SIGTERM...")
			pCancel()
		}		
	}()
	
	conn, err := grpc.Dial(_greeterServiceAddr, grpc.WithInsecure())
	if err != nil {
		calmwuUtils.Fatalf("grpc dial to %s failed. err:%s", err.Error())
	}
	defer conn.Close()

	client := protoHelloworld.NewGreeterClient(conn)
	name := _defaultName

	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	
	tickerCall := time.NewTicker(5 * time.Second)
	defer tickerCall.Stop()

L:
	for{
		select {
		case <-tickerCall.C:
			ctx, cancel := context.WithTimeout(pCtx, _defaultGRPCTimeout)
			defer cancel()
		
			resp, err := client.SayHello(ctx, &protoHelloworld.HelloRequest{
				Name: name,})
			if err != nil {
				calmwuUtils.Errorf("call greet.SayHello failed. err:%s", err.Error())
				break L
			} else {
				calmwuUtils.Debugf("index:%d call greet.SayHello resp:%#v", index++, resp)
			}
		case <-pCtx.Done():
			calmwuUtils.Debug("istio-simplegrpc-client recv exit notify")
			break L
		}
	}
}