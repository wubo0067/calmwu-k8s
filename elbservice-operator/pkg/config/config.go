/*
 * @Author: calmwu
 * @Date: 2020-05-30 22:33:25
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-05-31 00:01:07
 */

// config elbservice配置读取
package config

import (
	"io/ioutil"

	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ELBOpers 操作信息
type ELBOpers struct {
	BindURL   string `json:"BindURL" mapstructure:"BindURL"`
	UnBindURL string `json:"UnBindURL" mapstructure:"UnBindURL"`
}

// ELBEnvConfig 不同环境配置
type ELBEnvConfig struct {
	Env      string   `json:"Env" mapstructure:"Env"`
	ELBOpers ELBOpers `json:"ELBOpers" mapstructure:"ELBOpers"`
}

// ELBServiceConfig elbservice的配置信息
type ELBServiceConfig struct {
	CurrEnv       string         `json:"CurrEnv" mapstructure:"CurrEnv"`
	ELBEnvConfigs []ELBEnvConfig `json:"ELBEnvConfigs" mapstructure:"ELBEnvConfigs"`
}

var log = logf.Log.WithName("config")
var ELBServiceCfg = &ELBServiceConfig{}

func Init() error {
	err := unmarshalCfgFile("/etc/elbservice/elbservice_config.json")
	if err != nil {
		return err
	}
	return nil
}

func unmarshalCfgFile(file string) error {
	appFs := afero.NewOsFs()
	cfgFile, err := appFs.Open(file)
	if err != nil {
		log.Error(err, "Open file failed.", "file", file)
		return err
	}

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		log.Error(err, "read file failed.", "file", file)
		return err
	}

	err = ELBServiceCfg.UnmarshalJSON(cfgData)
	if err != nil {
		log.Error(err, "rELBServiceCfg.UnmarshalJSON failed.")
		return err
	}

	log.Info("ELBService_controller", "ELBServiceConfig", litter.Sdump(ELBServiceCfg))
	return nil
}
