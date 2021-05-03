/*
 * @Author: CALM.WU
 * @Date: 2021-04-29 10:57:39
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-29 18:07:56
 */

package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"pod.injector.nginx/pkg"
)

var (
	// Version 版本
	Version string
	// BuildTime 时间
	BuildTime string
)

func main() {
	var svrParameters pkg.SvrParamenters

	flag.IntVar(&svrParameters.Port, "port", 8443, "Webhook server port.")
	flag.StringVar(&svrParameters.CertFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&svrParameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&svrParameters.SidecarCfgFile, "sidecarCfgFile", "/etc/webhook/config/sidecarconfig.yaml", "File containing the mutation configuration.")
	flag.Parse()

	glog.Info("nginx-injector-pod-webhook-server starting...")

	pkg.LoadConfig(svrParameters.SidecarCfgFile)

	router := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		glog.Infof("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.POST("/inject", gin.WrapF(pkg.HandleInject))
	router.POST("/inject/", gin.WrapF(pkg.HandleInject))

	router.RunTLS(fmt.Sprintf("0.0.0.0:%d", svrParameters.Port), svrParameters.CertFile, svrParameters.KeyFile)
}
