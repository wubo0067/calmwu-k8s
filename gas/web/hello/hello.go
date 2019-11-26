/*
 * @Author: calm.wu
 * @Date: 2019-06-25 10:29:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 19:02:31
 */

package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sp_proto "gas/api/protobuf/srv/stringprocess"
	user_proto "gas/api/protobuf/srv/usermgr"
	hello_proto "gas/api/protobuf/web/hello"
	"gas/internal/utils/tracer"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"

	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/net/trace"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

type APIHello struct {
	Client client.Client
}

func (ah *APIHello) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Receive Hello.Call request")

	tr := trace.New("eci.vi.api.hello", "APIHello.Call")
	defer tr.Finish()

	// context
	ctx = trace.NewContext(ctx, tr)
	// 从context中获取值
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.Metadata{}
	}

	// add a unqiue request id to context
	if traceID, err := uuid.NewV4(); err == nil {
		tmd := metadata.Metadata{}
		for k, v := range md {
			tmd[k] = v
		}
		tmd["traceID"] = traceID.String()
		tmd["fromName"] = "eci.vi.api.hello"
		ctx = metadata.NewContext(ctx, tmd)
	}

	if req.Method == "GET" {
		log.Printf("req.Method is GET")

		name, ok := req.Get["name"]
		if !ok || len(name.Values) == 0 {
			return errors.BadRequest("eci.v1.api.hello", "no content")
		}

		spClient := sp_proto.NewStringProcessService("eci.v1.svr.stringprocess", ah.Client)

		res, err := spClient.ToUpper(ctx, &sp_proto.OriginalStrReq{OriginalString: name.Values[0]})
		if err != nil {
			log.Printf("invoke service:eci.v1.svr.stringprocess.ToUpper failed! reason:%s", err.Error())
			return errors.BadRequest("eci.v1.svr.stringprocess", err.Error())
		}

		// 在切割下
		spSplitClient := sp_proto.NewSplitProcessService("eci.v1.svr.stringprocess", ah.Client)
		splitRes, err := spSplitClient.Split(ctx, &sp_proto.OriginalStrReq{OriginalString: res.UpperString})
		if err != nil {
			log.Printf("invoke service:eci.v1.svr.stringprocess.Split failed! reason:%s", err.Error())
			return errors.BadRequest("eci.v1.svr.stringprocess", err.Error())
		}

		rsp.StatusCode = 200
		rsp.Body = fmt.Sprintf("[GET] Hello client %s!", splitRes.SplitStrs[0])
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

			userClient := user_proto.NewUserService("", ah.Client)
			stream, err := userClient.GetUserInfoServerStream(ctx, &user_proto.UserRequest{ID: 999})
			if err != nil {
				return errors.BadRequest("eci.v1.svr.user", err.Error())
			}

			for {
				rsp, err := stream.Recv()
				if err == io.EOF {
					// 服务端会关闭流，这里会收到EOF
					log.Printf("rsp recv completed!")
					break
				}

				if err != nil {
					log.Printf("recv err:%s\n", err.Error())
					return errors.BadRequest("eci.v1.svr.user", err.Error())
				}

				log.Printf("recv rsp:%#v\n", rsp)
				nameValue = fmt.Sprintf("%s-%d", rsp.Name, rsp.Age)
			}
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

	t, io, err := tracer.NewTracer("eci.v1.api.hello", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}

	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 本身是个服务
	service := micro.NewService(
		micro.Name("eci.v1.api.hello"),
		micro.RegisterTTL(time.Second*30), // 这里会设置consul.go Start中的ttl参数，func (s *rpcServer) Register() error {
		micro.RegisterInterval(time.Second*10),
		micro.Context(ctx),                                                        // 停服控制，用信号处理
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())), // 调链跟踪
	)

	service.Init()

	// graceful shutdown
	err = service.Server().Init(
		server.Wait(nil),
	)
	if err != nil {

		log.Fatal(err.Error())
	}

	hello_proto.RegisterHelloHandler(service.Server(), &APIHello{service.Client()})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
	return
}
