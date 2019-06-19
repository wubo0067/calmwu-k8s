/*
 * @Author: calm.wu 
 * @Date: 2019-06-18 10:17:00 
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-18 11:06:32
 */

 package main

 import (
	"fmt"
	"context"
	proto_hello "hello-microsvr/api/protobuf/hello"

	micro "github.com/micro/go-micro"
 )

 func main() {
	service := micro.NewService(micro.Name("hello.client"))
	service.Init()

	helloService := proto_hello.NewHelloService("Hello", service.Client())
	res, err := helloService.Ping(context.TODO(), &proto_hello.Request{Name: "World ^-^"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Msg)
 }