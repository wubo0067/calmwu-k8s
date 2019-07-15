/*
 * @Author: calm.wu
 * @Date: 2019-07-09 14:29:32
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 15:51:21
 */

package gateway

// https://github.com/Allenxuxu/microservices/blob/master/micro/main.go

import (
	"gas/internal/utils/tracer"
	"log"

	"github.com/Allenxuxu/microservices/lib/wrapper/tracer/opentracing/stdhttp"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	opentracing "github.com/opentracing/opentracing-go"
)

func init() {
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
}

func Main() {
	t, io, err := tracer.NewTracer("eci.v1.gateway", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}

	defer io.Close()
	opentracing.SetGlobalTracer(t)

	cmd.Init()
}
