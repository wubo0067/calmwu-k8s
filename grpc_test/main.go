/*
 * @Author: calm.wu
 * @Date: 2019-06-12 16:36:23
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-12 18:56:53
 */

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	user "grpc_test/user"
)

type UserService struct {
	// 实现 User 服务的业务对象
}

// UserService 实现了 User 服务接口中声明的所有方法
func (userService *UserService) UserIndex(ctx context.Context, in *user.UserIndexRequest) (*user.UserIndexResponse, error) {
	log.Printf("receive user index request: page %d page_size %d", in.Page, in.PageSize)

	return &user.UserIndexResponse{
		Err: 0,
		Msg: "success",
		Data: []*user.UserEntity{
			{Name: "big_cat", Age: 28},
			{Name: "sqrt_cat", Age: 29},
		},
	}, nil
}

func (userService *UserService) UserView(ctx context.Context, in *user.UserViewRequest) (*user.UserViewResponse, error) {
	log.Printf("receive user view request: uid %d", in.Uid)

	return &user.UserViewResponse{
		Err:  0,
		Msg:  "success",
		Data: &user.UserEntity{Name: "james", Age: 28},
	}, nil
}

func (userService *UserService) UserPost(ctx context.Context, in *user.UserPostRequest) (*user.UserPostResponse, error) {
	log.Printf("receive user post request: name %s password %s age %d", in.Name, in.Password, in.Age)

	return &user.UserPostResponse{
		Err: 0,
		Msg: "success",
	}, nil
}

func (userService *UserService) UserDelete(ctx context.Context, in *user.UserDeleteRequest) (*user.UserDeleteResponse, error) {
	log.Printf("receive user delete request: uid %d", in.Uid)

	return &user.UserDeleteResponse{
		Err: 0,
		Msg: "success",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:5556")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	user.RegisterUserServer(grpcServer, &UserService{})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return
}
