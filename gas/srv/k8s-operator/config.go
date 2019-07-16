/*
 * @Author: calm.wu
 * @Date: 2019-07-15 14:26:02
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-15 16:59:16
 */

package k8soperator

import (
	"fmt"
	"sync"

	"gas/internal/utils/config"

	"github.com/pquerna/ffjson/ffjson"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

const (
	// consul配置路径
	srvConfigPath         = "config/srv"
	k8sOperatorConfigNode = "k8soperator"
)

// 具体的配置信息
type srvK8sOperatorConfigInfo struct {
	LogPath       string `json:"logPath" mapstructure:"logPath"`             // log的存放路径
	JaegerSvrAddr string `json:"jaegerSvrAddr" mapstructure:"jaegerSvrAddr"` // 调链跟踪服务的地址信息ip:port
}

// 服务的配置模块
type srvK8sOperatorConfigMgr struct {
	configInfo   *srvK8sOperatorConfigInfo
	consulConfig *config.ConfigMgr
	initOnce     sync.Once
	wg           sync.WaitGroup
	quitCh       chan struct{}
	monitor      sync.Mutex
}

func newSrvK8sOperatorConfigMgr() *srvK8sOperatorConfigMgr {
	return &srvK8sOperatorConfigMgr{
		configInfo:   new(srvK8sOperatorConfigInfo),
		consulConfig: nil,
		quitCh:       make(chan struct{}),
	}
}

func (cm *srvK8sOperatorConfigMgr) init() error {
	var err error
	// ensure init once
	cm.initOnce.Do(func() {
		cm.consulConfig, err = config.NewConfigMgr(srvConfigPath)
		if err != nil {
			return
		}

		// 启动watch goroutine
		cm.wg.Add(1)
		go cm.watchCfg()
	})

	return err
}

func (cm *srvK8sOperatorConfigMgr) stop() {
	close(cm.quitCh)
	cm.wg.Wait()
	cm.consulConfig.Stop()
}

func (cm *srvK8sOperatorConfigMgr) watchCfg() {
	defer func() {
		cm.wg.Done()
		if err := recover(); err != nil {
			stackInfo := calm_utils.CallStack(1)
			fmt.Println(stackInfo)
		}
		fmt.Println("srvK8sOperatorConfigMgr config watch routine exit")
	}()

L:
	for {
		select {
		case _, ok := <-cm.consulConfig.ChangeNtfChan:
			if ok {
				// 更新配置对象
				cm.getConfigInfo()
			} else {
				fmt.Printf("config change notify channel is closed!\n")
				break L
			}
		case <-cm.quitCh:
			break L
		}
	}
	return
}

func (cm *srvK8sOperatorConfigMgr) getConfigInfo() (*srvK8sOperatorConfigInfo, error) {
	cm.monitor.Lock()
	defer cm.monitor.Unlock()

	configData := cm.consulConfig.GetConfigData(k8sOperatorConfigNode)
	//fmt.Printf("configData:%s\n", string(configData))
	// unmarshal
	cm.configInfo = new(srvK8sOperatorConfigInfo)
	err := ffjson.Unmarshal(configData, cm.configInfo)
	if err != nil {
		//fmt.Printf("err:%s\n", err.Error())
		return nil, err
	}
	//fmt.Printf("configInfo:%v\n", cm.configInfo)
	return cm.configInfo, nil
}
