/*
 * @Author: calm.wu
 * @Date: 2019-11-26 10:42:24
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-11-26 14:50:26
 */

package strupper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	proto_upper "multi-call/api/protobuf/srv-upper"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	logger *log.Logger
)

type StrUpperImpl struct{}

func (sui *StrUpperImpl) ToUpper(ctx context.Context, in *proto_upper.StrUpperReq, out *proto_upper.StrUpperRes) error {
	logger.Println("------------receive ToUpper call------------")

	ctxDeadlineTime, ctxDeadlineOK := ctx.Deadline()
	logger.Printf("ctx deadline time:%s ok:%v\n", ctxDeadlineTime.String(), ctxDeadlineOK)

	out.OriginalString = strings.ToUpper(in.OriginalString)
	logger.Println("ToUpper call completed")
	return nil
}

func Main() {
	logger = calm_utils.NewSimpleLog(nil)

	logger.Printf("srv-strupper start\n")

	svrReg := consul.NewRegistry(registryOptions)

	// 这里传入各种option函数，用于修改server的options
	service := micro.NewService(
		micro.Name("sci.v1.svr.strupper"),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(svrReg),
	)

	service.Init()

	proto_upper.RegisterStrUpperProcessHandler(service.Server(), &StrUpperImpl{})

	logger.Printf("micro register RegisterStrUpperProcessHandler\n")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
	// 设定tcp检测时间间隔
	//ops.Context = context.WithValue(context.Background(), "consul_tcp_check", 5*time.Second)
}
