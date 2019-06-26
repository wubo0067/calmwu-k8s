/*
 * @Author: calm.wu
 * @Date: 2019-06-25 10:29:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-06-25 10:30:12
 */

package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	hello_proto "gas/api/protobuf/web/hello"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
)

type APIHello struct{}

func (ah *APIHello) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Receive Hello.Call request")

	if req.Method == "GET" {
		name, ok := req.Get["name"]
		if !ok || len(name.Values) == 0 {
			return errors.BadRequest("eci.v1.api.hello", "no content")
		}

		rsp.StatusCode = 200
		rsp.Body = fmt.Sprintf("[GET] Hello client %s!", name.GetValues()[0])
		return nil
	} else if req.Method == "POST" {
		ct, ok := req.Header["Content-Type"]
		if !ok || len(ct.Values) == 0 {
			return errors.BadRequest("go.micro.api.hello", "need content-type")
		}

		if ct.Values[0] != "application/json" {
			return errors.BadRequest("go.micro.api.hello", "expect application/json")
		}

		// parse body
		var body map[string]interface{}
		json.Unmarshal([]byte(req.Body), &body)
		if nameValueI, exists := body["name"]; exists {
			nameValue := nameValueI.(string)
			rsp.StatusCode = 200
			rsp.Body = fmt.Sprintf("[POST] Hello client %s!", nameValue)
			return nil
		}
	}

	return errors.BadRequest("eci.v1.api.hello", "method is invalid!")
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			sig := <-sigCh
			switch sig {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				cancel()
				return
			}
		}
	}()
}

func Main() {
	// 停服控制
	ctx, cancel := context.WithCancel(context.Background())
	setupSignalHandler(cancel)

	// 本身是个服务
	service := micro.NewService(
		micro.Name("eci.v1.api.hello"),
		micro.RegisterTTL(time.Second*30), // 这里会设置consul.go Start中的ttl参数，func (s *rpcServer) Register() error {
		micro.RegisterInterval(time.Second*10),
		micro.Context(ctx), // 停服控制，用信号处理
	)

	service.Init()

	// graceful shutdown
	err := service.Server().Init(
		server.Wait(true),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	hello_proto.RegisterHelloHandler(service.Server(), new(APIHello))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	return
}
