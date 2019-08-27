/*
 * @Author: calm.wu
 * @Date: 2019-08-27 16:27:45
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-27 17:16:02
 */

package config

import (
	"io/ioutil"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// NSPData 配置数据
type NSPData struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

// K8SClusterData 多集群配置数据
type K8SClusterData struct {
	K8SClusterID string `json:"k8sclusterid" mapstructure:"k8sclusterid"`
	KubeCfg      string `json:"kubecfg" mapstructure:"kubecfg"`
}

// SrvIPResMgrConfigData 服务的配置数据
type SrvIPResMgrConfigData struct {
	NSPData           NSPData          `json:"nsp" mapstructure:"nsp"`
	K8SClusterDataLst []K8SClusterData `json:"k8sclustser_list" mapstructure:"k8sclustser_list"`
}

var (
	configData *SrvIPResMgrConfigData
	configFile string
	guard      sync.Mutex
)

// LoadConfig 加载配置数据
func LoadConfig(configFile string) error {
	err := calm_utils.PathExist(configFile)
	if err != nil {
		os.Exit(-1)
	}
	configFile = configFile

	configData = new(SrvIPResMgrConfigData)

	cfgFile, err := os.Open(configFile)
	if err != nil {
		return errors.Wrapf(err, "open %s failed", configFile)
	}
	defer cfgFile.Close()

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return errors.Wrapf(err, "read %s failed", configFile)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(cfgData, configData)
	if err != nil {
		return errors.Wrap(err, "json Unmarshal failed.")
	}

	calm_utils.ZLog.Debugf("ipresmgr-svr config:%+v", configData)
	return nil
}

// ReloadConfig 重新加载配置
func ReloadConfig() {
	cfgFile, err := os.Open(configFile)
	if err != nil {
		calm_utils.ZLog.Errorf("open config file %s failed.", configFile)
		return
	}
	defer cfgFile.Close()

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		calm_utils.ZLog.Errorf("read config file %s failed.", configFile)
		return
	}

	newConfigData := new(SrvIPResMgrConfigData)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(cfgData, newConfigData)
	if err != nil {
		calm_utils.ZLog.Error("json umarshal config data failed.")
		return
	}

	guard.Lock()
	defer guard.Unlock()
	configData = newConfigData
	newConfigData = nil

	calm_utils.ZLog.Infof("reload config %s successed", configFile)
}
