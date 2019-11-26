/*
 * @Author: calm.wu
 * @Date: 2019-11-26 11:20:18
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-11-26 14:47:54
 */

package strproc

import (
	"context"
	"fmt"
	"time"

	proto_split "multi-call/api/protobuf/srv-split"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

func Main() {
	// 定义服务，可以传入其它可选参数
	svrReg := consul.NewRegistry(registryOptions)
	service := micro.NewService(
		micro.Name("strproc.client"),
		micro.Registry(svrReg))
	service.Init()

	// 创建新的客户端
	strprocClient := proto_split.NewStrSplitProcessService("", service.Client())

	// 调用greeter
	// 这个context没有设置超时时间，callOptions修改请求超时时间
	ctx := context.Background()

	rsp, err := strprocClient.Split(ctx, &proto_split.StrSplitReq{
		OriginalString: "Hello-world",
	}, func(op *client.CallOptions) {
		op.RequestTimeout = 10 * time.Second
	})
	if err != nil {
		fmt.Println(err)
	}

	// 打印响应请求
	fmt.Println(rsp.SplitStrs)
}

func registryOptions(ops *registry.Options) {
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
	// 设定tcp检测时间间隔
	//ops.Context = context.WithValue(context.Background(), "consul_tcp_check", 5*time.Second)
}
