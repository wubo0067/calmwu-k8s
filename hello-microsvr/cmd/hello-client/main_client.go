/*
 * @Author: calm.wu
 * @Date: 2019-06-18 10:17:00
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-18 11:06:32
 */

package main

import (
	"context"
	"fmt"
	proto_hello "hello-microsvr/api/protobuf/hello"
	"time"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

func main() {
	svrReg := consul.NewRegistry(registryOptions)

	service := micro.NewService(micro.Name("hello.client"), micro.Registry(svrReg),)
	service.Init()

	helloService := proto_hello.NewHelloService("Hello", service.Client())
	res, err := helloService.Ping(context.TODO(), &proto_hello.Request{Name: "World ^-^"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Msg)
}

func registryOptions(ops *registry.Options) {
	ops.Timeout = time.Second * 5
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
}
