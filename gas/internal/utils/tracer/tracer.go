/*
 * @Author: calm.wu
 * @Date: 2019-07-09 11:00:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 11:21:09
 */

package tracer

import (
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

//
func NewTracer(serviceName string, jaegerAddr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	// set Jaeger report receive address
	sender, err := jaeger.NewUDPTransport(jaegerAddr, 0)
	if err != nil {
		return nil, nil, err
	}

	// create jaeger reporter
	reporter := jaeger.NewRemoteReporter(sender)
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}
