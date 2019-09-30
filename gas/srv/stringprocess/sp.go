/*
 * @Author: calm.wu
 * @Date: 2019-06-26 14:37:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-23 16:36:13
 */

package stringprocess

import (
	"context"
	"log"
	"strings"
	"time"

	//"github.com/micro/cli"
	sp_proto "gas/api/protobuf/srv/stringprocess"
	"gas/internal/utils/tracer"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"golang.org/x/net/trace"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

type StringProcessHandlerImpl struct{}

func (spi *StringProcessHandlerImpl) ToUpper(ctx context.Context, in *sp_proto.OriginalStrReq, out *sp_proto.UpperStrRes) error {
	log.Println("service: eci.v1.svr.stringprocess handler:StringProcessHandler method:ToUpper")

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

type SplitProcessHandlerImpl struct{}

func (sph *SplitProcessHandlerImpl) Split(ctx context.Context, in *sp_proto.OriginalStrReq, out *sp_proto.SplitStrRes) error {
	log.Println("service: eci.v1.svr.stringprocess handler:SplitProcessHandler method:Split")

	out.SplitStrs = strings.Split(in.OriginalString, "-")
	return nil
}

// logWrapper is a handler wrapper
// 传入实际的接口函数
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[wrapper] server request: %v", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}

type StringProcessServiceHandler struct {
	sp_proto.StringProcessHandler
	sp_proto.SplitProcessHandler
}

func Main() {
	//svrTransport := grpc.NewTransport(transportOptions)
	t, io, err := tracer.NewTracer("eci.v1.svr.stringprocess", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}

	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 这里传入各种option函数，用于修改server的options
	service := micro.NewService(
		// 这个名字必须是protobuf的service名字
		// 这里是有namespace的
		micro.Name("eci.v1.svr.stringprocess"),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		//micro.Transport(svrTransport),
		micro.WrapHandler(logWrapper), // 这里是handlerwapper，是对回调方法的封装
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())), // 调链跟踪，NewHandlerWrapper返回一个HandlerWrapper函数对象
	)

	// 服务初始化
	service.Init()

	// 这里会注册多个handler
	serviceHandler := &StringProcessServiceHandler{
		&StringProcessHandlerImpl{},
		&SplitProcessHandlerImpl{},
	}

	sp_proto.RegisterStringProcessHandler(service.Server(), serviceHandler)
	sp_proto.RegisterSplitProcessHandler(service.Server(), serviceHandler)

	// sp_proto.RegisterStringProcessHandler(service.Server(), new(StringProcessHandlerImpl))
	// sp_proto.RegisterSplitProcessHandler(service.Server(), new(SplitProcessHandlerImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func transportOptions(ops *transport.Options) {
}
