/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:45:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-13 16:17:19
 */

package srv

import (
	"fmt"
	"os"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/k8s"
	"pci-ipresmgr/pkg/ipresmgr/nsp"
	"pci-ipresmgr/pkg/ipresmgr/storage"
	"pci-ipresmgr/pkg/ipresmgr/storage/mysql"

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

var (
	storeMgr storage.StoreMgr
)

func initLog(logFilePath string, srvInstID string) {
	err := calm_utils.CheckDir(logFilePath)
	if err != nil {
		os.Exit(-1)
	}

	calm_utils.InitDefaultZapLog(fmt.Sprintf("%s/%s.log", logFilePath, srvInstID),
		zapcore.DebugLevel, 1)
}

// func setupSignalHandler(cancel context.CancelFunc) {
// 	sigCh := make(chan os.Signal, 1)
// 	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
// 	go func() {
// 		for {
// 			sig := <-sigCh
// 			switch sig {
// 			case syscall.SIGINT:
// 				fallthrough
// 			case syscall.SIGTERM:
// 				calm_utils.Info("-------------catch shutdown signal-------------")
// 				cancel()
// 				return
// 			case syscall.SIGUSR1:
// 				config.ReloadConfig()
// 				k8s.DefaultK8SClient.LoadMultiClusterClient(config.GetK8SClusterCfgDataLst())
// 			case syscall.SIGUSR2:
// 				calm_utils.DumpStacks()
// 			}
// 		}
// 	}()
// }

// SvrMain 服务的入口
func SvrMain(c *cli.Context) error {
	// 获取参数
	configFile := c.String("conf")
	logFilePath := c.String("logpath")
	srvInstID := fmt.Sprintf("ipresmgr-svr_%d", c.Int("id"))
	listenAddr := c.String("ip")
	listenPort := c.Int("port")

	initLog(logFilePath, srvInstID)

	calm_utils.Infof("-------------%s start running-------------", srvInstID)

	err := config.LoadConfig(configFile)
	if err != nil {
		calm_utils.Fatalf("LoadConfig %s failed, err:%s", configFile, err.Error())
	}

	loadOk := k8s.DefaultK8SClient.LoadMultiClusterClient(config.GetK8SClusterCfgDataLst())
	if !loadOk {
		calm_utils.Fatal("LoadMultiClusterClient failed")
	}

	// 初始化存储
	storeMgr = mysql.NewMysqlStoreMgr()
	err = storeMgr.Start(func(opts *storage.StoreOptions) {
		storeCfgData := config.GetStoreCfgData()
		opts.SrvInstID = srvInstID
		opts.StoreSvrAddr = storeCfgData.MysqlAddr
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

	// 初始化nsp
	nsp.NSPInit(config.GetNspServerAddr())

	// 启动服务
	startWebSrv(listenAddr, listenPort)

	calm_utils.Info("ipresmgr-srv http server exiting")
	return nil
}
