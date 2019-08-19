/*
 * @Author: calm.wu
 * @Date: 2019-07-09 14:29:32
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 15:51:21
 */

package gateway

// https://github.com/Allenxuxu/microservices/blob/master/micro/main.go

import (
	"gas/internal/utils/tracer"
	"log"
	"time"

	"github.com/Allenxuxu/microservices/lib/wrapper/tracer/opentracing/stdhttp"
	go_micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	opentracing "github.com/opentracing/opentracing-go"
)

func init() {
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
}

func Main() {
	t, io, err := tracer.NewTracer("eci.v1.gateway", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}

	svrReg := consul.NewRegistry(registryOptions)

	defer io.Close()
	opentracing.SetGlobalTracer(t)

	cmd.Init(
		go_micro.RegisterTTL(time.Second*5),
		go_micro.RegisterInterval(time.Second*2),
		go_micro.Registry(svrReg),
	)
}

func registryOptions(ops *registry.Options) {
	//ops.Timeout = time.Second * 5
	ops.Addrs = []string{"127.0.0.1:8500"}
	//ops.Context = context.WithValue(context.TODO(), "consul_tcp_check", time.Duration(5*time.Second))
}
