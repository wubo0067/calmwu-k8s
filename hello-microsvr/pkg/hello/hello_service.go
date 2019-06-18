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

	"hello-microsvr/api/protobuf/hello"

	micro "github.com/micro/go-micro"
)

type HelloService struct{}

func (h *HelloService) Ping(ctx context.Context, req *hello.Request, res *hello.Response) error {
	res.Msg = "Hello" + req.Name
	return nil
}

func Main() {
	fmt.Println("hello.Main")

	service := micro.NewService(micro.Name("Hello"))
	service.Init()

	hello.RegisterHelloHandler(service.Server(), new(HelloService))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
