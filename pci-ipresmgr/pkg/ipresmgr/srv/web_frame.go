/*
 * @Author: calm.wu
 * @Date: 2019-08-27 17:32:13
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 12:16:07
 */

package srv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	httpSrv *http.Server
)

func registerHandler(router *gin.Engine) {
	// webhook接口
	wbV1Group := router.Group("/v1/ippool")
	wbV1Group.POST("/create", wbCreateIPPool)
	wbV1Group.POST("/release", wbReleaseIPPool)
	wbV1Group.POST("/scale", wbScaleIPPool)

	// cni接口
	cniV1Group := router.Group("/v1/ip")
	cniV1Group.POST("/require", cniRequireIP)
	cniV1Group.POST("/release", cniReleaseIP)

	// maintain接口
	maintainGroup := router.Group("/v1/maintain")
	maintainGroup.POST("/unbindIP", maintainUnbindIP)
	maintainGroup.POST("/releaseIP", maintainReleaseIP)
}

func startWebSrv(listenAddr string, listenPort int) error {
	gin.SetMode(gin.DebugMode)
	ginRouter := gin.New()
	ginRouter.Use(calm_utils.GinLogger())
	ginRouter.Use(calm_utils.GinRecovery())

	// 注册健康检查接口
	ginRouter.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 注册pprof接口
	ginpprof.Wrap(ginRouter)

	// 注册业务接口
	registerHandler(ginRouter)

	httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", listenAddr, listenPort),
		Handler: ginRouter,
	}

	// 启动监听
	go func() {
		calm_utils.Infof("ipresmgr-svr listen on %s:%d", listenAddr, listenPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			calm_utils.Fatalf("Listen %s:%d failed. err:%s", listenAddr, listenPort, err.Error())
		}
	}()
	return nil
}

func shutdownWebSrv() {
	if httpSrv != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 给shutdown 5秒时间
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		}
		select {
		case <-shutdownCtx.Done():
			calm_utils.Info("delay 5 seconds for graceful shutdown")
		}
		calm_utils.Info("ipresmgr-srv http server exiting")
	}
	return
}
