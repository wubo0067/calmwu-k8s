/*
 * @Author: calm.wu
 * @Date: 2019-11-26 10:42:24
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-11-26 14:50:15
 */

package strsplit

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	proto_split "multi-call/api/protobuf/srv-split"
	proto_upper "multi-call/api/protobuf/srv-upper"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	logger *log.Logger
)

type StrSplitImpl struct {
	client client.Client
}

func (ssi *StrSplitImpl) Split(ctx context.Context, in *proto_split.StrSplitReq, out *proto_split.StrSplitRes) error {
	logger.Println("------------receive Split call------------")

	ctxDeadlineTime, ctxDeadlineOK := ctx.Deadline()
	logger.Printf("ctx deadline time:%s ok:%v\n", ctxDeadlineTime.String(), ctxDeadlineOK)

	//time.Sleep(5 * time.Second)
	// 调用toupper
	// 这样就限制了本次调用的时间为3秒
	ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()
	strUpperClient := proto_upper.NewStrUpperProcessService("", ssi.client)
	upperRsp, err := strUpperClient.ToUpper(ctx, &proto_upper.StrUpperReq{
		OriginalString: in.GetOriginalString(),
	}, func(op *client.CallOptions) {
		op.RequestTimeout = 7 * time.Second
	})

	if err != nil {
		err = errors.Wrap(err, "strUpperClient.ToUpper call failed.")
		logger.Println(err.Error())
	}

	out.SplitStrs = strings.Split(upperRsp.OriginalString, "-")
	logger.Println("Split call completed")
	return nil
}

func Main() {
	logger = calm_utils.NewSimpleLog(nil)

	logger.Printf("srv-strsplit start\n")

	svrReg := consul.NewRegistry(registryOptions)

	// 这里传入各种option函数，用于修改server的options
	service := micro.NewService(
		micro.Name("sci.v1.svr.strsplit"),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(svrReg),
	)

	service.Init()

	proto_split.RegisterStrSplitProcessHandler(service.Server(), &StrSplitImpl{client: service.Client()})

	logger.Printf("micro register RegisterStrSplitProcessHandler\n")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	ops.Addrs = []string{fmt.Sprintf("%s:%d", "127.0.0.1", 8500)}
	// 设定tcp检测时间间隔
	//ops.Context = context.WithValue(context.Background(), "consul_tcp_check", 5*time.Second)
}
