/*
 * @Author: calm.wu
 * @Date: 2019-08-27 17:32:13
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2019-11-23 14:07:14
 */

package srv

import (
	"fmt"
	"net/http"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/k8s"
	"syscall"

	"github.com/DeanThompson/ginpprof"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// var (
// 	httpSrv *http.Server
// )

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
	maintainGroup.POST("/unbindip", maintainForceUnbindIP)
	maintainGroup.POST("/release/ippool", maintainForceReleaseK8SResourceIPPool)
	maintainGroup.POST("/release/podip", maintainForceReleasePodIP)
}

func preHookSigUsr1Reload() {
	config.ReloadConfig()
	k8s.DefaultK8SClient.LoadMultiClusterClient(config.GetK8SClusterCfgDataLst())
}

func preHookSigUsr2DumpStack() {
	calm_utils.DumpStacks()
}

func preHookSigIntShutdown() {
	calm_utils.Warnf("ipresmgr-srv receive SIGINT for shutdown")
}

func preHookSigTermShutdown() {
	calm_utils.Warnf("ipresmgr-srv receive SIGTERM for shutdown")
}

func preHookSigHupRestart() {
	calm_utils.Warnf("ipresmgr-srv receive SIGHUP for restart")
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

	httpSrv := endless.NewServer(fmt.Sprintf("%s:%d", listenAddr, listenPort), ginRouter)

	// 注册新号响应回调
	httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
		httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
		preHookSigIntShutdown)
	httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1] = append(
		httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1],
		preHookSigUsr1Reload)
	httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR2] = append(
		httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR2],
		preHookSigUsr2DumpStack)
	httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGTERM] = append(
		httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGTERM],
		preHookSigTermShutdown)
	httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGHUP] = append(
		httpSrv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGHUP],
		preHookSigHupRestart)

	// 启动监听
	calm_utils.Infof("ipresmgr-svr listen on %s:%d", listenAddr, listenPort)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		calm_utils.Errorf("Listen %s:%d failed. err:%s", listenAddr, listenPort, err.Error())
		return err
	}

	return nil
}
