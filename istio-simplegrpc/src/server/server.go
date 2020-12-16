/*
 * @Author: CALM.WU
 * @Date: 2020-11-29 12:07:36
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-11-29 22:04:22
 */

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	protoHelloworld "istio-simplegrpc/proto/helloworld"
	protoPerson "istio-simplegrpc/proto/person"

	"github.com/sanity-io/litter"
	calmwuUtils "github.com/wubo0067/calmwu-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

//
type IstioSimpleGRPCServerImpl struct {
	// 这里必须嵌入，https://github.com/grpc/grpc-go/issues/3669
	protoHelloworld.UnimplementedGreeterServer
	protoPerson.UnimplementedPersonRegistryServer
}

var (
	_index    = 0
	_hostName = ""
)

const (
	_defaultCallType = "CallType"
)

// SayHello 测试接口
func (isgsi *IstioSimpleGRPCServerImpl) SayHello(ctx context.Context, in *protoHelloworld.HelloRequest) (*protoHelloworld.HelloReply, error) {
	_index++

	// 获取http header
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(_defaultCallType)
	calmwuUtils.Debugf("index:%d Greeter.SayHello called, name: %s, CallType:%#v", _index, in.Name, vals)

	return &protoHelloworld.HelloReply{
		Message: fmt.Sprintf("srv-host:%s index:%d Hello %s", _hostName, _index, in.Name),
	}, nil
}

// CreateReservation 测试接口
func (isgsi *IstioSimpleGRPCServerImpl) CreateReservation(ctx context.Context, in *protoHelloworld.CreateReservationRequest) (*protoHelloworld.Reservation, error) {
	_index++

	// 获取http header
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(_defaultCallType)
	calmwuUtils.Debugf("index:%d Greeter.CreateReservation called, Reservation: %s, CallType:%#v", _index, litter.Sdump(in.Reservation), vals)

	return &protoHelloworld.Reservation{
		Id:        strconv.FormatInt(int64(_index), 10),
		Title:     fmt.Sprintf("%s-%d", in.Reservation.Title, _index),
		Venue:     in.Reservation.Venue,
		Room:      in.Reservation.Room,
		Timestamp: in.Reservation.Timestamp,
		Attendees: in.Reservation.Attendees,
	}, nil
}

func (isgsi *IstioSimpleGRPCServerImpl) EchoTimeout(ctx context.Context, in *protoHelloworld.EchoRequest) (*protoHelloworld.EchoReply, error) {
	_index++

	// 获取http header
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(_defaultCallType)
	calmwuUtils.Debugf("index:%d Greeter.EchoTimeout called, name: %s, CallType:%#v", _index, in.Message, vals)

	// 故意延迟5秒钟，测试istio的超时和重试功能
	time.Sleep(5 * time.Second)

	return &protoHelloworld.EchoReply{
		Message: fmt.Sprintf("srv-host:%s index:%d Echo %s", _hostName, _index, in.Message),
	}, nil
}

func (isgsi *IstioSimpleGRPCServerImpl) Lookup(ctx context.Context, in *protoPerson.Person) (*protoPerson.Person, error) {
	_index++

	calmwuUtils.Debugf("index:%d Greeter.Lookup called, Person: %s", _index, litter.Sdump(in))
	return &protoPerson.Person{
		Name: in.Name,
		Age:  in.Age,
		Addr: &protoPerson.Address{
			HouseNum:   fmt.Sprintf("Lookup-HouseNum-%d", _index),
			Building:   fmt.Sprintf("Lookup-Building-%d", _index),
			Street:     fmt.Sprintf("Lookup-Street-%d", _index),
			Locality:   fmt.Sprintf("Lookup-Locality-%d", _index),
			City:       fmt.Sprintf("Lookup-City-%d", _index),
			PostalCode: fmt.Sprintf("Lookup-PostalCode-%d", _index),
		},
	}, nil
}

func (isgsi *IstioSimpleGRPCServerImpl) Create(ctx context.Context, in *protoPerson.Person) (*protoPerson.Person, error) {
	_index++

	calmwuUtils.Debugf("index:%d Greeter.Create called, Person: %s", _index, litter.Sdump(in))

	return &protoPerson.Person{
		Name: in.Name,
		Age:  in.Age,
		Addr: &protoPerson.Address{
			HouseNum:   fmt.Sprintf("Create-HouseNum-%d", _index),
			Building:   fmt.Sprintf("Create-Building-%d", _index),
			Street:     fmt.Sprintf("Create-Street-%d", _index),
			Locality:   fmt.Sprintf("Create-Locality-%d", _index),
			City:       fmt.Sprintf("Create-City-%d", _index),
			PostalCode: fmt.Sprintf("Create-PostalCode-%d", _index),
		},
	}, nil
}

var (
	simpleGrpcSrvImpl *IstioSimpleGRPCServerImpl       = &IstioSimpleGRPCServerImpl{}
	_                 protoHelloworld.GreeterServer    = simpleGrpcSrvImpl
	_                 protoPerson.PersonRegistryServer = simpleGrpcSrvImpl
)

func main() {
	calmwuUtils.Debug("istio-simplegrpc-server now start.")

	_hostName = os.Getenv("HOSTNAME")

	listen, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		calmwuUtils.Fatalf("failed to listen: %v", err.Error())
	}

	grpcSrv := grpc.NewServer()
	// 注册服务
	protoHelloworld.RegisterGreeterServer(grpcSrv, simpleGrpcSrvImpl)
	protoPerson.RegisterPersonRegistryServer(grpcSrv, simpleGrpcSrvImpl)

	reflection.Register(grpcSrv)
	if err := grpcSrv.Serve(listen); err != nil {
		calmwuUtils.Fatal("failed to serve: %v", err.Error())
	}
}
