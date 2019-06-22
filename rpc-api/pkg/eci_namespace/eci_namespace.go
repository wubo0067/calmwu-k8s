/*
 * @Author: calm.wu
 * @Date: 2019-06-21 16:25:25
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-21 16:32:36
 */

package eci_namespace

import (
	"context"
	"fmt"
	"log"

	ns_proto "rpc-api/api/protobuf/eci-namespace"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

type NamespaceService struct{}

func (ns *NamespaceService) GetNamespace(ctx context.Context, in *ns_proto.CallRequest, out *ns_proto.CallResponse) error {
	log.Println("receive getnamespace request")

	out.NamespaceInfo = in.Name + "_namespace"
	return nil
}

func Main() {
	fmt.Println("eci-namespace service")

	service := micro.NewService(
		// 这个名字必须是protobuf的service名字
		micro.Name("eci.v1.api.NamespaceSvr"),
	)

	// 一定要加上init，不会解析输入参数，
	service.Init(
		micro.Action(func(c *cli.Context) {
			fmt.Println("eci-namespace Service Init")
		}),
	)

	ns_proto.RegisterNamespaceSvrHandler(service.Server(), new(NamespaceService))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
