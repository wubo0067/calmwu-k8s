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
	_ignoreNamespaces = hashset.New()

	_defaultIgnoreNamespaces = []string{"kube-system", "kube-public", "nginx-injector-pod-webhook"}
)

// LoadConfig read config data from configmap
func LoadConfig(configFile string) (*SidecarConfig, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = errors.Wrapf(err, "read config file:%s failed.", configFile)
		glog.Error(err.Error())

		return nil, err
	}

	var cfg SidecarConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		err = errors.Wrapf(err, "yaml decode config file:%s failed.", configFile)
		glog.Error(err.Error())

		return nil, err
	}

	// 读取环境变量
	ignoreNamespacesStr := os.Getenv("IGNORE_NAMESPACES")
	nsList := strings.Split(ignoreNamespacesStr, ":")
	nsList = append(nsList, _defaultIgnoreNamespaces...)

	for _, ignoreNS := range nsList {
		_ignoreNamespaces.Add(ignoreNS)
	}

	glog.Infof("Sidecar config:\n%s", string(data))
	glog.Infof("ENV IGNORE_NAMESPACES: %s", _ignoreNamespaces.Values())

	return &cfg, nil
}

func isIgnoreNamespace(ns string) bool {
	return _ignoreNamespaces.Contains(ns)
}
