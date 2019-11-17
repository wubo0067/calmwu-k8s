/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:21:48
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-27 16:37:58
 */

// Package ipresmgr for
package ipresmgr

import (
	"fmt"
	"os"
	"pci-ipresmgr/pkg/ipresmgr/srv"

	"github.com/micro/cli"
	_ "go.uber.org/automaxprocs"
)

var (
	appname = "ipresmgr-svr"
)

// Main 地址资源管理服务的入口
func Main(buildTime string, version string) {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version:%s Buildtime:%s\n", version, buildTime)
	}

	app := cli.NewApp()
	app.Name = appname
	app.Usage = "Management of the fixed ip of the container"
	app.Action = srv.SvrMain
	app.Flags = srv.SvrFlags

	app.Run(os.Args)
	fmt.Println(appname, "exit!")
}
