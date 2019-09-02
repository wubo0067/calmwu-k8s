/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:45:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-30 14:32:21
 */

package srv

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/store"
	"pci-ipresmgr/pkg/ipresmgr/store/mysql"
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

func initLog(logFilePath string, srvInstID string) {
	err := calm_utils.CheckDir(logFilePath)
	if err != nil {
		os.Exit(-1)
	}

	calm_utils.InitDefaultZapLog(fmt.Sprintf("%s/%s.log", logFilePath, srvInstID),
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
				calm_utils.Info("-------------catch shutdown signal-------------")
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
	srvInstID := fmt.Sprintf("ipresmgr-svr_%d", c.Int("id"))
	listenAddr := c.String("ip")
	listenPort := c.Int("port")

	initLog(logFilePath, srvInstID)

	err := config.LoadConfig(configFile)
	if err != nil {
		calm_utils.Fatalf("LoadConfig %s failed, err:%s", configFile, err.Error())
	}

	// 信号控制
	ctx, cancel := context.WithCancel(context.Background())
	setupSignalHandler(cancel)

	// 初始化存储
	storeMgr := mysql.NewMysqlStoreMgr()
	err = storeMgr.Start(ctx, func(opts *store.StoreOptions) {
		storeCfgData := config.GetStoreCfgData()
		opts.SrvInstID = srvInstID
		opts.Addr = storeCfgData.MysqlAddr
		opts.User = storeCfgData.User
		opts.Passwd = storeCfgData.Passwd
		opts.DBName = storeCfgData.DBName
		opts.IdelConnectCount = storeCfgData.IdelConnectCount
		opts.MaxOpenConnectCount = storeCfgData.MaxOpenConnectCount
		opts.ConnectMaxLifeTime = storeCfgData.ConnectMaxLifeTime
	})
	if err != nil {
		calm_utils.Fatalf("storeMgr start failed, err:%s", err.Error())
	}
	defer storeMgr.Stop()

	err = storeMgr.Register(listenAddr, listenPort)
	if err != nil {
		calm_utils.Errorf("register self failed. err:%s", err.Error())
		storeMgr.Stop()
		return err
	}
	defer storeMgr.UnRegister()

	// 初始化web
	err = startWebSrv(listenAddr, listenPort)
	if err != nil {
		return err
	}

	// 等待退出信号
	select {
	case <-ctx.Done():
		calm_utils.Info("ipresmgr-svr will shutdown")
	}

	// 退出清理

	// 停止web服务
	shutdownWebSrv()
	// 停止存储
	//storeMgr.UnRegister(fmt.Sprintf("ipresmgr-svr_%d", srvInstID))
	//storeMgr.Stop()
	return nil
}
