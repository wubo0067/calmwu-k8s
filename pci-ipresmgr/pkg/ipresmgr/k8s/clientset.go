/*
 * @Author: calm.wu
 * @Date: 2019-09-13 10:49:47
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2019-11-23 15:14:46
 */
// Package k8s for
package k8s

import (
	"time"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func kubeCfgLoadGetter(content []byte) clientcmd.KubeconfigGetter {
	return func() (*api.Config, error) {
		return clientcmd.Load(content)
	}
}

// NewClientSetByKubeCfgContent 根据kubcfg内容构造Clientset对象
func NewClientSetByKubeCfgContent(content []byte) (*kubernetes.Clientset, error) {
	if len(content) == 0 {
		return nil, errors.New("Content args is invalid")
	}

	conf, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeCfgLoadGetter(content))
	if err != nil {
		err = errors.Wrap(err, "clientcmd BuildConfigFromKubeconfigGetter failed.")
		calm_utils.Error(err)
		return nil, err
	}
	// 设置请求超时时间
	conf.Timeout = time.Duration(5 * time.Second)

	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		err = errors.Wrap(err, "kubernetes NewForConfig failed.")
		calm_utils.Error(err)
		return nil, err
	}

	return clientSet, nil
}
