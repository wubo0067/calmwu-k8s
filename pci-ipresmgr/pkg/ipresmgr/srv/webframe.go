/*
 * @Author: calm.wu
 * @Date: 2019-08-27 17:32:13
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-27 19:21:39
 */

package svr

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	httpSrv *http.Server
)

func ginRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := calm_utils.CallStack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				calm_utils.ZLog.Errorf("[Recovery] panic recovered:\n%s\n%s\n%s", calm_utils.Bytes2String(httprequest), err, stack)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		calm_utils.ZLog.Debugf("%s latency:%s", c.Request.RequestURI, latency.String())
	}
}

func startWebSrv(listenAddr string, listenPort int) error {
	gin.SetMode(gin.DebugMode)
	ginRouter := gin.New()
	ginRouter.Use(ginLogger())
	ginRouter.Use(ginRecovery())

	httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", listenAddr, listenPort),
		Handler: ginRouter,
	}

	// 启动监听
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			calm_utils.ZLog.Fatalf("Listen %s:%d failed. err:%s", listenAddr, listenPort, err.Error())
		}
	}()
	return nil
}

func shutdownWebSrv() {
	// 有个超时保证
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
	}
	select {
	case <-shutdownCtx.Done():
		calm_utils.ZLog.Info("")
	}
	return
}
