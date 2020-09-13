/*
 * @Author: calmwu
 * @Date: 2020-09-13 20:47:12
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-09-13 20:56:24
 */

package client

import (
	"testing"

	"kube-utils/vendor/k8s.io/client-go/discovery"

	"k8s.io/client-go/tools/clientcmd"
)

func TestDiscoveryClient(t *testing.T) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		t.Fatalf("build restCfg failed. err:%s", err.Error())
	}

	// discovery.NewDiscoveryClientForConfigg函数通过config实例化discoveryClient对象
	// 用于发现kube-apiserver所支持的资源组(Group)，资源版本(Versions)，资源信息(Resources)
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restCfg)
	if err != nil {
		t.Fatalf("create discovery client failed. err:%s", err.Error())
	}

	//
	APIGroupList, APIResourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		t.Fatalf("discovery client ServerGroupsAndResources failed. err:%s", err.Error())
	}

	for index := range APIGroupList {
		t.Logf("[%d] apiGroup:%s", index, APIGroupList[index].String())
	}

	for index := range APIResourceList {
		t.Logf("[%d] apiResource:%s", index, APIResourceList[index].String())
	}
}
