/*
 * @Author: calm.wu
 * @Date: 2019-08-27 16:27:45
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-13 15:09:33
 */

// Package config ....
package config

import (
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// NSPData 配置数据
type NSPData struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

// K8SClusterCfgData 多集群配置数据
type K8SClusterCfgData struct {
	K8SClusterID string `json:"k8sclusterid" mapstructure:"k8sclusterid"`
	KubeCfg      string `json:"kubecfg" mapstructure:"kubecfg"`
}

// StoreCfgData 数据库配置
type StoreCfgData struct {
	MysqlAddr           string `json:"mysqladdr" mapstructure:"mysqladdr"`
	User                string `json:"user" mapstructure:"user"`
	Passwd              string `json:"passwd" mapstructure:"passwd"`
	DBName              string `json:"dbname" mapstructure:"dbname"`
	IdelConnectCount    int    `json:"idlconnects" mapstructure:"idlconnects"`
	MaxOpenConnectCount int    `json:"maxopenconnects" mapstructure:"maxopenconnects"`
	ConnectMaxLifeTime  string `json:"connectmaxlifetime" mapstructure:"connectmaxlifetime"`
}

// SrvIPResMgrConfigData 服务的配置数据
type SrvIPResMgrConfigData struct {
	NSPData                        NSPData             `json:"nsp" mapstructure:"nsp"`
	K8SClusterCfgDataLst           []K8SClusterCfgData `json:"k8sclustser_list" mapstructure:"k8sclustser_list"`
	StoreData                      StoreCfgData        `json:"store" mapstructure:"store"`
	K8SResourceAddrLeasePeriodSecs int                 `json:"k8sResourceAddrLeasePeriodSecs" mapstructure:"k8sResourceAddrLeasePeriodSecs"`
}

var (
	configFileName string
	configVal      atomic.Value
)

// LoadConfig 加载配置数据
func LoadConfig(configFile string) error {
	err := calm_utils.PathExist(configFile)
	if err != nil {
		os.Exit(-1)
	}
	configFileName = configFile

	configData := new(SrvIPResMgrConfigData)

	cfgFile, err := os.Open(configFileName)
	if err != nil {
		return errors.Wrapf(err, "open:[%s] failed", configFileName)
	}
	defer cfgFile.Close()

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return errors.Wrapf(err, "read:[%s] failed", configFileName)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(cfgData, configData)
	if err != nil {
		return errors.Wrap(err, "json Unmarshal failed.")
	}

	configVal.Store(configData)
	calm_utils.Debugf("ipresmgr-svr config:%s", litter.Sdump(configData))
	return nil
}

// ReloadConfig 重新加载配置
func ReloadConfig() {
	cfgFile, err := os.Open(configFileName)
	if err != nil {
		calm_utils.Errorf("open:[%s] failed.", configFileName)
		return
	}
	defer cfgFile.Close()

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		calm_utils.Errorf("read:[%s] failed.", configFileName)
		return
	}

	newConfigData := new(SrvIPResMgrConfigData)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(cfgData, newConfigData)
	if err != nil {
		calm_utils.Error("json umarshal config data failed.")
		return
	}

	configVal.Store(newConfigData)

	calm_utils.Infof("reload config:[%s] successed. newConfigData:%s", configFileName, litter.Sdump(newConfigData))
}

// GetNspServerAddr 获取nsp服务的地址
func GetNspServerAddr() string {
	configData := configVal.Load().(*SrvIPResMgrConfigData)
	nspSrvAddr := configData.NSPData.Addr
	return nspSrvAddr
}

// GetStoreCfgData 获取配置
func GetStoreCfgData() StoreCfgData {
	configData := configVal.Load().(*SrvIPResMgrConfigData)
	return configData.StoreData
}

// GetK8SResourceAddrLeasePeriodSecs 获得地址租期时间，单位秒
func GetK8SResourceAddrLeasePeriodSecs() time.Duration {
	configData := configVal.Load().(*SrvIPResMgrConfigData)
	return time.Duration(configData.K8SResourceAddrLeasePeriodSecs)
}

// GetK8SClusterCfgDataLst 得到多集群的kubecfg配置信息
func GetK8SClusterCfgDataLst() []K8SClusterCfgData {
	configData := configVal.Load().(*SrvIPResMgrConfigData)
	return configData.K8SClusterCfgDataLst
}
