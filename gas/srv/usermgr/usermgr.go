/*
 * @Author: calm.wu
 * @Date: 2019-06-30 11:34:01
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-30 16:30:57
 */

package usermgr

import (
	"context"
	"io"
	"log"
	"time"

	user_proto "gas/api/protobuf/srv/usermgr"

	"github.com/micro/go-micro"
)

type UserServiceImpl struct{}

var users = map[int32]user_proto.UserResponse{
	1: {Name: "Ken Thompson", Age: 75},
	2: {Name: "Jeff Deam", Age: 40},
	3: {Name: "John Camark", Age: 43},
}

// Server side stream
func (usi *UserServiceImpl) GetUserInfoServerStream(ctx context.Context, req *user_proto.UserRequest,
	stream user_proto.UserService_GetUserInfoServerStreamStream) error {
	log.Printf("Receive req:%#v Server side stream\n", req)

	// server返回多行数据
	for _, user := range users {
		stream.Send(&user)
	}

	stream.Close()
	return nil
}

// Bidirectinal stream
func (usi *UserServiceImpl) GetUserInfoBidirectionalStream(ctx context.Context,
	stream user_proto.UserService_GetUserInfoBidirectionalStreamStream) error {
	log.Println("Receive Bidirectional stream\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client close stream")
			break
		}

		if err != nil {
			log.Printf("stream Recv err:%s\n", err.Error())
			return err
		}

		log.Printf("Receive req:%#v\n", req)

		if err = stream.Send(&users[req.ID%3]); err != nil {
			log.Printf("stream Send err:%s\n", err.Error())
			return err
		}
	}
	return nil
}

func Main() {
	service := micro.NewService(
		micro.Name("eci.v1.server.user"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*3),
	)

	service.Init()

	user_proto.RegisterUserServiceHandler(service.Server(), new(UserServiceImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
