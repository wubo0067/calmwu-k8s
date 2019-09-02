/*
 * @Author: calm.wu
 * @Date: 2019-08-27 16:27:45
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-01 10:35:45
 */

package config

import (
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"

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
	NSPData           NSPData          `json:"nsp" mapstructure:"nsp"`
	K8SClusterDataLst []K8SClusterData `json:"k8sclustser_list" mapstructure:"k8sclustser_list"`
	StoreData         StoreCfgData     `json:"store" mapstructure:"store"`
}

var (
	configFileName string
	guard          sync.Mutex
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
	calm_utils.ZLog.Debugf("ipresmgr-svr config:%+v", configData)
	return nil
}

// ReloadConfig 重新加载配置
func ReloadConfig() {
	cfgFile, err := os.Open(configFileName)
	if err != nil {
		calm_utils.ZLog.Errorf("open:[%s] failed.", configFileName)
		return
	}
	defer cfgFile.Close()

	cfgData, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		calm_utils.ZLog.Errorf("read:[%s] failed.", configFileName)
		return
	}

	newConfigData := new(SrvIPResMgrConfigData)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(cfgData, newConfigData)
	if err != nil {
		calm_utils.ZLog.Error("json umarshal config data failed.")
		return
	}

	configVal.Store(newConfigData)

	calm_utils.ZLog.Infof("reload config:[%s] successed", configFileName)
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
