/*
 * @Author: calm.wu
 * @Date: 2019-06-26 14:37:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-26 14:46:19
 */

package stringprocess

import (
	"log"
	"context"
	"strings"

	//"github.com/micro/cli"
	sp_proto "gas/api/protobuf/srv/stringprocess"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
)

type StringProcessImpl struct{}

func (spi *StringProcessImpl) ToUpper(ctx context.Context, in *sp_proto.OriginalStrReq, out *sp_proto.UpperStrRes) error {
	log.Println("receive eci.v1.svr.stringprocess.ToUpper request")

	out.UpperString = strings.ToUpper(in.OriginalString)
	return nil
}

func Main() {
	svrTransport := grpc.NewTransport(transportOptions)

	service := grpc.NewService(
		// 这个名字必须是protobuf的service名字
		// 这里是有namespace的
		micro.Name("eci.v1.svr.stringprocess"),
		micro.Transport(svrTransport)
	)

	service.Init()

	sp_proto.RegisterStringProcessHandler(service.Server(), new(StringProcessImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}	
}

func transportOptions(ops *transport.Options) {
}
