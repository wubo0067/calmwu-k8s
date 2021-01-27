/*
 * @Author: CALM.WU
 * @Date: 2021-01-14 11:34:49
 * @Last Modified by:   CALM.WU
 * @Last Modified time: 2021-01-14 11:34:49
 */

package main

import (
	"os"
	"time"

	"ng-kubeportforward/portforward"

	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildOutOfClusterConfig(kubeCfgPath string) (*rest.Config, error) {
	kubeconfigPath := kubeCfgPath
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func main() {
	restCfg, _ := buildOutOfClusterConfig("/root/.kube/config")
	stopCh := make(chan struct{})

	localPort, _ := portforward.GetFreePort()
	pfPod := portforward.NewPortForward(restCfg, "kata-ngdp-8568bd4758-65l69", "default", localPort, 80, stopCh)
	err := pfPod.Start()
	if err != nil {
		calmUtils.Error(err.Error())
		return
	}

	time.Sleep(3 * time.Second)
	pfPod.Stop()
	time.Sleep(time.Second)
}
