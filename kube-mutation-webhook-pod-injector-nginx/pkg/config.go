/*
 * @Author: CALM.WU
 * @Date: 2021-04-29 14:26:23
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-29 14:28:40
 */

// Package pkg is implement nginx injector to pod
package pkg

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

// SvrParamenters is server parameters
type SvrParamenters struct {
	Port           int    // webhook server port
	CertFile       string // path to the x509 certificate for https
	KeyFile        string // path to the x509 private key matching `CertFile`
	SidecarCfgFile string // path to sidecar injector configuration file
}

// SidecarConfig is inject container info
type SidecarConfig struct {
	Containers []corev1.Container `yaml:"containers"`
	Volumes    []corev1.Volume    `yaml:"volumes"`
}

var (
	_ignoreNamespaces        = hashset.New()
	_defaultIgnoreNamespaces = []string{"kube-system", "kube-public", "nginx-injector-pod-webhook"}
	_configFile              string
	_configWatcher           Watcher
	_sidecarConfig           *SidecarConfig
	// StopWatchCh stop watch file
	StopWatchCh chan struct{}
)

// LoadConfig read config data from configmap
func LoadConfig(configFile string) error {
	if _configWatcher == nil {
		var err error
		_configWatcher, err = NewFileWatcher(configFile)
		if err != nil {
			err = errors.Wrap(err, "Create file watcher failed.")
			glog.Infof(err.Error())
			return err
		}

		_configFile = configFile
		_configWatcher.SetUpdateNotify(loadConfig)

		StopWatchCh = make(chan struct{})
		go _configWatcher.Run(StopWatchCh)
	}

	return loadConfig(configFile)
}

func loadConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = errors.Wrapf(err, "read config file:%s failed.", configFile)
		glog.Error(err.Error())
		return err
	}

	if _sidecarConfig != nil {
		_sidecarConfig = nil
	}

	_sidecarConfig = new(SidecarConfig)
	if err := yaml.Unmarshal(data, _sidecarConfig); err != nil {
		err = errors.Wrapf(err, "yaml decode config file:%s failed.", configFile)
		glog.Error(err.Error())
		return err
	}

	// 读取环境变量
	ignoreNamespacesStr := os.Getenv("IGNORE_NAMESPACES")
	nsList := strings.Split(ignoreNamespacesStr, ":")
	nsList = append(nsList, _defaultIgnoreNamespaces...)

	if _ignoreNamespaces.Size() != 0 {
		_ignoreNamespaces.Clear()
	}

	for _, ignoreNS := range nsList {
		_ignoreNamespaces.Add(ignoreNS)
	}

	glog.Infof("Sidecar config:\n%s", string(data))
	glog.Infof("ENV IGNORE_NAMESPACES: %s", _ignoreNamespaces.Values())

	return nil
}

func isIgnoreNamespace(ns string) bool {
	return _ignoreNamespaces.Contains(ns)
}

func getSidecarConfig() *SidecarConfig {
	return _sidecarConfig
}
