/*
 * @Author: calm.wu
 * @Date: 2019-06-26 14:37:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-30 11:34:51
 */

package stringprocess

import (
	"context"
	"log"
	"strings"

	//"github.com/micro/cli"
	sp_proto "gas/api/protobuf/srv/stringprocess"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"golang.org/x/net/trace"
)

type StringProcessImpl struct{}

func (spi *StringProcessImpl) ToUpper(ctx context.Context, in *sp_proto.OriginalStrReq, out *sp_proto.UpperStrRes) error {
	log.Println("receive eci.v1.svr.stringprocess.ToUpper request")

	// 从ctx中获取metadata数据
	md, _ := metadata.FromContext(ctx)
	traceID := md["traceID"]
	fromName := md["fromName"]

	// 获取tr
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("fromName: %s traceID %s", fromName, traceID)
	} else {
		log.Println("from context get trace object is nil")
	}

	out.UpperString = strings.ToUpper(in.OriginalString)
	return nil
}

// logWrapper is a handler wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[wrapper] server request: %v", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}

func Main() {
	//svrTransport := grpc.NewTransport(transportOptions)

	service := micro.NewService(
		// 这个名字必须是protobuf的service名字
		// 这里是有namespace的
		micro.Name("eci.v1.svr.stringprocess"),
		//micro.Transport(svrTransport),
		micro.WrapHandler(logWrapper), // 这里是handlerwapper，是对回调方法的封装
	)

	service.Init()

	sp_proto.RegisterStringProcessHandler(service.Server(), new(StringProcessImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func transportOptions(ops *transport.Options) {
}
