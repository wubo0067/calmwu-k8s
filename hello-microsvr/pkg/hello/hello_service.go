/*
 * @Author: calm.wu
 * @Date: 2019-06-17 15:23:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-23 16:37:47
 */

package hello

import (
	"context"
	"fmt"
	"time"

	proto_hello "hello-microsvr/api/protobuf/hello"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	codec_proto "github.com/micro/go-micro/codec/proto"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
)

type HelloService struct{}

func (h *HelloService) Ping(ctx context.Context, req *proto_hello.Request, res *proto_hello.Response) error {
	res.Msg = "Hello" + req.Name
	return nil
}

func Main() {
	fmt.Println("hello.Main")

	// 自定义server
	//svr := server.NewServer(serverOptions)

	// 构建一个grpc transport，--transport=grpc是使用默认的配置
	svrTransport := grpc.NewTransport(transportOptions)

	// 使用consul，这用自定义的配置了，而不是consul的默认配置
	// 返回一个registry.Registry interface，这个对象会赋值给micro.Options.Registry成员，
	// 赋值前这个对象已经根据配置文件填充了配置信息
	svrReg := consul.NewRegistry(registryOptions)

	// 这样写是取代默认的options中的组件
	service := micro.NewService(
		micro.Name("Hello"),           // 这个是服务名字
		micro.Registry(svrReg),        // 返回一个函数指针，为service.Options.Registry和service.Options.Server.Registry进行赋值
		micro.Transport(svrTransport), // 同样如此，这样取代默认的组件，用自己生成的
	)

	// 这里负责服务的初始化工作, 例如存储、传输等等
	// https://github.com/micro-in-cn/tutorials/tree/master/microservice-in-micro/part1
	service.Init(
		micro.Action(func(c *cli.Context) {
			fmt.Println("Hello Service Init")
		}),
	)

	// 其实核心是Server，默认是rpcServer，在这里可以做一定的修改对server
	// 这样写是对默认的server组件进行初始化，设置特定的参数
	service.Server().Init(serverOptions)

	proto_hello.RegisterHelloHandler(service.Server(), new(HelloService))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// 用来初始化consul的option
func registryOptions(ops *registry.Options) {
	ops.Timeout = time.Second * 5
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
	// 设定tcp检测时间间隔
	ops.Context = context.WithValue(context.Background(), "consul_tcp_check", 5*time.Second)
}

/*
这个是默认的配置for transport
github.com/micro/go-micro/transport.Options {
	Addrs: []string len: 0, cap: 0, nil,
	Codec: github.com/micro/go-micro/codec.Marshaler nil,
	Secure: false,
	TLSConfig: *crypto/tls.Config nil,
	Timeout: 0,
	Context: context.Context nil,}

*/
func transportOptions(ops *transport.Options) {
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 10008)}
	ops.Codec = codec_proto.Marshaler{}
}

func serverOptions(ops *server.Options) {
	fmt.Printf("server Options:%#v\n", ops)
	ops.Address = fmt.Sprintf("%s:%d", "127.0.0.1", 10009)
	ops.Id = "10008"
}
