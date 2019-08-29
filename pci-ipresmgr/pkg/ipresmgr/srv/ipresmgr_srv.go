/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:45:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-29 19:24:53
 */

package srv

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/store"
	"syscall"

	"github.com/micro/cli"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"go.uber.org/zap/zapcore"
)

var (
	// SvrFlags 命令参数
	SvrFlags = []cli.Flag{
		cli.IntFlag{
			Name:  "id",
			Value: 1,
			Usage: "Server instance ID",
		},
		cli.StringFlag{
			Name:  "logpath",
			Value: "log",
			Usage: "The path to the log file",
		},
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "ipresmgr server config file path",
		},
		cli.StringFlag{
			Name:  "ip, i",
			Value: "0.0.0.0",
			Usage: "ipresmgr server listen ip",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 9000,
			Usage: "ipresmgr server listen port",
		},
	}
)

func initLog(logFilePath string, srvInstID int) {
	err := calm_utils.CheckDir(logFilePath)
	if err != nil {
		os.Exit(-1)
	}

	calm_utils.InitDefaultZapLog(fmt.Sprintf("%s/ipresmgr-svr_%d.log", logFilePath, srvInstID),
		zapcore.DebugLevel, 1)
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigCh
			switch sig {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				calm_utils.Info("catch shutdown signal")
				cancel()
				return
			case syscall.SIGUSR1:
				config.ReloadConfig()
			case syscall.SIGUSR2:
				calm_utils.DumpStacks()
			}
		}
	}()
}

// SvrMain 服务的入口
func SvrMain(c *cli.Context) error {
	// 获取参数
	configFile := c.String("conf")
	logFilePath := c.String("logpath")
	srvInstID := c.Int("id")
	listenAddr := c.String("ip")
	listenPort := c.Int("port")

	initLog(logFilePath, srvInstID)

	err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("LoadConfig %s failed, err:%s", configFile, err.Error())
	}

	// 信号控制
	ctx, cancel := context.WithCancel(context.Background())
	setupSignalHandler(cancel)

	// 初始化存储
	backEndStore := store.NewStore()
	err = backEndStore.Start(func(opts *store.StoreOptions) {
	})

	// 初始化web
	startWebSrv(listenAddr, listenPort)

	// 等待退出信号
	select {
	case <-ctx.Done():
		calm_utils.Info("ipresmgr-svr will shutdown")
	}

	// 退出清理

	// 停止web服务
	shutdownWebSrv()
	// 停止存储
	backEndStore.Stop()
	return nil
}
