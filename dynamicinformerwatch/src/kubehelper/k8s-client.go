/*
 * @Author: CALM.WU
 * @Date: 2021-04-14 10:29:47
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-14 17:27:44
 */

// Package kubehelper is tools for k8s
package kubehelper

import (
	"github.com/pkg/errors"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetRestConfig 得到kubecfg对象
func GetRestConfig(kubeCfgPath string) (*rest.Config, error) {
	// if kubeCfgPath == "" {
	// 	kubeCfgPath = os.Getenv("KUBECONFIG")
	// 	if kubeCfgPath == "" {
	// 		kubeCfgPath = os.Getenv("HOME") + "/.kube/config"
	// 	}
	// }

	var conf *rest.Config
	var err error

	if kubeCfgPath == "" {
		conf, err = rest.InClusterConfig()
		if err != nil {
			err = errors.Wrap(err, "in cluster")
			calmUtils.Error(err.Error())
			return nil, err
		}
	} else {
		conf, err = clientcmd.BuildConfigFromFlags("", kubeCfgPath)
		if err != nil {
			err = errors.Wrap(err, "out cluster")
			calmUtils.Error(err.Error())
			return nil, err
		}
	}

	return conf, nil
}

// MakeClientSet returns a k8s clientset to the request from inside of cluster
func MakeClientSet(kubeCfgPath string) kubernetes.Interface {
	restCfg, err := GetRestConfig(kubeCfgPath)
	if err != nil {
		calmUtils.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		calmUtils.Fatalf("could not generate kubernetes client for config, err: %s", err.Error())
	}

	return clientset
}

func MakeDynamicClient(kubeCfgPath string) dynamic.Interface {
	restCfg, err := GetRestConfig(kubeCfgPath)
	if err != nil {
		calmUtils.Fatal(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		calmUtils.Fatalf("could not generate dynamic client for config, err: %s", err.Error())
	}

	return dynamicClient
}
