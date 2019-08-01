/*
 * @Author: calm.wu
 * @Date: 2019-07-03 10:50:28
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-04 11:14:14
 */

package config

import (
	"log"
	"sync"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/consul"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// https://github.com/micro/go-micro/tree/master/config/source/consul
// https://github.com/micro-in-cn/tutorials/tree/master/microservice-in-micro/part4

type ConfigMgr struct {
	microConfig   config.Config
	watcher       config.Watcher
	wait          sync.WaitGroup
	configPath    string
	ChangeNtfChan chan struct{}
}

func NewConfigMgr(configPath string) (*ConfigMgr, error) {
	configMgr := &ConfigMgr{
		microConfig:   config.NewConfig(),
		watcher:       nil,
		configPath:    configPath,
		ChangeNtfChan: make(chan struct{}),
	}

	var err error

	// 如果使用toml
	//enc := toml.NewEncoder()

	consulConfSource := consul.NewSource(
		consul.WithAddress("127.0.0.1:8500"),
		consul.WithPrefix(configPath),
		consul.StripPrefix(true),
		//source.WithEncoder(enc), 指定编解码类型
	)

	// 设置source
	err = configMgr.microConfig.Load(consulConfSource)
	if err != nil {
		return nil, err
	}

	configMgr.microConfig.Sync()

	// 启动watch
	configMgr.watcher, err = configMgr.microConfig.Watch()
	if err != nil {
		return nil, err
	}

	configMgr.wait.Add(1)
	go configMgr.watchConsulConfig()

	return configMgr, nil
}

func (cm *ConfigMgr) watchConsulConfig() {
	defer func() {
		cm.wait.Done()
		if err := recover(); err != nil {
			log.Printf("recover error:%s\n", err)
		}
	}()

	for {
		// 监听consul配置变化
		cs, err := cm.watcher.Next()
		if err != nil {
			// watch stop
			break
		}
		log.Printf("%s content change, %s", cm.configPath, calm_utils.Bytes2String(cs.Bytes()))

		select {
		// 一定要加default，防止routine泄漏
		case cm.ChangeNtfChan <- struct{}{}:
		default:
		}
	}
}

// 停止
func (cm *ConfigMgr) Stop() {
	close(cm.ChangeNtfChan)
	cm.watcher.Stop()
	cm.wait.Wait()
}

// 获取配置数据，默认是json数据，需要反序列化
func (cm *ConfigMgr) GetConfigData(path ...string) []byte {
	value := cm.microConfig.Get(path...)
	return value.Bytes()
}
