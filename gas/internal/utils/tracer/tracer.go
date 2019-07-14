/*
 * @Author: calm.wu
 * @Date: 2019-07-09 11:00:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 11:21:09
 */

package tracer

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

const ginTracerKey = "ginTracerKey-context"

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

//
func GinTracerWrapper(c *gin.Context) {
	md := make(map[string]string)
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
	defer sp.Finish()

	if err := opentracing.GlobalTracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	c.Set(ginTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	}
}

// GinContextWithSpan 返回context
func GinContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(ginTracerKey)
	if exist == false {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}

/*
	router := gin.Default()
	r := router.Group("/user")
	r.Use(tracer.GinTracerWrapper)

	在方法中调用GinContextWithSpan

	ctx, ok := tracer.GinContextWithSpan(c)，将这个ctx传递给rpc的的方法

	https://github.com/Allenxuxu/microservices/blob/master/api/user/
*/
