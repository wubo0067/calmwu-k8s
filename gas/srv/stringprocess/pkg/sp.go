/*
 * @Author: calm.wu
 * @Date: 2019-06-26 14:37:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-26 14:46:19
 */

package stringprocess

import (
	//"github.com/micro/cli"
	sp_proto "gas/svr/stringprocess/protobuf"
	"github.com/micro/go-micro"
)

type StringProcessImpl struct{}

func (spi *StringProcessImpl) ToUpper(ctx context.Context, in *sp_proto.OriginalStrReq, out *sp_proto.UpperStrRes) error {
	return nil
}

func Main() {
	service := micro.NewService(
		// 这个名字必须是protobuf的service名字
		// 这里是有namespace的
		micro.Name("eci.v1.svr.stringprocess"),
	)

	service.Init()

	sp_proto.RegisterStringProcessHandler(service.Server(), new(StringProcessImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}	
}
