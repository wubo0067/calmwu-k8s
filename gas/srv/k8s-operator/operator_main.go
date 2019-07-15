/*
 * @Author: calm.wu
 * @Date: 2019-07-15 14:18:50
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-15 17:05:29
 */

package k8soperator

import (
	"context"
	"fmt"
	"io"
	"time"
	"log"

	"gas/internal/utils/tracer"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/micro/cli"
	opentracing "github.com/opentracing/opentracing-go"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
)

const (
	srvName = "pci.v1.srv.k8s.operator"
)

var (
	srvConfig *srvK8sOperatorConfigMgr
)

// 服务实际的入口函数
func Main() {
	// 初始化配置
	err := initConfig()
	if err != nil {
		return
	}
	defer srvConfig.stop()
	// 得到配置信息
	configInfo, err := srvConfig.getConfigInfo()
	if err != nil {
		return
	}

	// 初始化log

	ctx, cancel := context.WithCancel(context.Background())
	// 初始化信号
	setupSignalHandler(cancel)

	// 注册consul，自定义参数
	svrReg := consul.NewRegistry(registryOptions)

	// 调链跟踪
	tracerIO, err := initTracer(configInfo.jaegerSvrAddr)
	if err != nil {
		return
	}
	defer tracerIO.Close()

	service := micro.NewService(
		micro.Name(srvName),
		micro.RegisterTTL(time.Second*30), // 这里会设置consul.go Start中的ttl参数，func (s *rpcServer) Register() error {
		micro.RegisterInterval(time.Second*10),
		micro.Registry(svrReg),
		micro.Context(ctx),                                                        // 停服控制，用信号处理
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())), // 调链跟踪
	)	

	service.Init(
		micro.Action(func(c *cli.Context) {
			fmt.Printf("%s Service Init", srvName)
		}),
	)

	service.Server().Init(
		server.Wait(nil), // graceful
	)	

	// 注册服务接口

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	return
}

// 初始化配置
func initConfig() error {
	srvConfig = newSrvK8sOperatorConfigMgr()
	err := srvConfig.init()
	if err != nil {
		return err
	}
	return nil
}

// 自定义consul注册参数
func registryOptions(ops *registry.Options) {
	ops.Timeout = time.Second * 5
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
	// 设定tcp检测时间间隔
	ops.Context = context.WithValue(context.Background(), "consul_tcp_check", 5*time.Second)
}

// 初始化调链跟踪
func initTracer(jaegerSvrAddr string) (io.Closer, error) {
	t, io, err := tracer.NewTracer(srvName, jaegerSvrAddr)
	if err != nil {
		return io, err
	}

	opentracing.SetGlobalTracer(t)
	return io, err
}
