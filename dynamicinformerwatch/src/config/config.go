/*
 * @Author: CALM.WU
 * @Date: 2021-04-14 11:19:07
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-14 17:10:47
 */

// Package config is load config data
package config

import (
	"sync/atomic"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
)

// watchK8SResource 监控资源对象
type watchK8SResource struct {
	ResGVK string `koanf:"name"`
}

// watchResConfig 配置对象
type watchResConfig struct {
	KubeCfg        string             `koanf:"kubecfg"`
	Namespace      string             `koanf:"namespace"`
	WatchResources []watchK8SResource `koanf:"resources"`
}

const (
	configFilePath = "../config/config.json"
)

var (
	kf         = koanf.New(".")
	configData atomic.Value
)

// LoadConf 加载配置文件
func LoadConf() error {
	// 读取配置文件
	cp := file.Provider(configFilePath)
	if err := kf.Load(cp, json.Parser()); err != nil {
		err = errors.Wrapf(err, "load file: %s", configFilePath)
		calmUtils.Error(err.Error())
		return err
	}

	// 反序列化
	cfg := &watchResConfig{}
	if err := kf.Unmarshal("watch-conf", &cfg); err != nil {
		err = errors.Wrapf(err, "unmarshal file: %s", configFilePath)
		calmUtils.Error(err.Error())
		return err
	}

	configData.Store(cfg)

	// 监控配置文件变化
	cp.Watch(func(event interface{}, err error) {
		if err != nil {
			calmUtils.Errorf("watch file: %s error, %s", err.Error())
			return
		}

		// 重新加载
		calmUtils.Debugf("config file: %s changed. reloading...", configFilePath)
		kf.Load(cp, json.Parser())
		kf.Unmarshal("watch-conf", &cfg)
		calmUtils.Debug(litter.Sdump(cfg))
		configData.Store(cfg)
	})

	return nil
}

// GetConfData 获取配置数据
func GetConfData() *watchResConfig {
	cfData := configData.Load().(*watchResConfig)
	return cfData
}
