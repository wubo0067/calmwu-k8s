/*
 * @Author: calm.wu
 * @Date: 2019-06-17 15:23:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-18 11:14:58
 */

package hello

import (
	"context"
	"fmt"

	proto_hello "hello-microsvr/api/protobuf/hello"

	micro "github.com/micro/go-micro"
	"github.com/micro/cli"
)

type HelloService struct{}

func (h *HelloService) Ping(ctx context.Context, req *proto_hello.Request, res *proto_hello.Response) error {
	res.Msg = "Hello" + req.Name
	return nil
}

func Main() {
	fmt.Println("hello.Main")

	service := micro.NewService(micro.Name("Hello"))

	// 这里负责服务的初始化工作, 例如存储、传输等等
	// https://github.com/micro-in-cn/tutorials/tree/master/microservice-in-micro/part1
	service.Init(
		micro.Action(func(c *cli.Context){
			fmt.Println("Hello Service Init")
		}),
	)

	proto_hello.RegisterHelloHandler(service.Server(), new(HelloService))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
