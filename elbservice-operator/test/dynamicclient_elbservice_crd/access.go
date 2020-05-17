/*
 * @Author: calmwu
 * @Date: 2020-05-16 16:23:52
 * @Last Modified by:   calmwu
 * @Last Modified time: 2020-05-16 16:23:52
 */

// 使用dynamic client去访问CRD资源
// https://soggy.space/namespaced-crds-dynamic-client/

package main

import (
	"flag"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	// 通过命令 kubectl api-resources 查看resource，使用kubectl api-versions查看gv信息
	_elbServiceGVR = schema.GroupVersionResource{
		Group:    "k8s.calmwu.org",
		Version:  "v1alpha1",
		Resource: "elbservices",
	}
)

func main() {
	calm_utils.Debug("Starting test")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//out cluster
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		calm_utils.Panic(err.Error())
	}

	// 构造dynamic client
	dynClient, errClient := dynamic.NewForConfig(config)
	if errClient != nil {
		calm_utils.Fatalf("Create dynamic client failed. err:%v", errClient)
	}

	elbServiceClient := dynClient.Resource(_elbServiceGVR)

	// 获取的是没有对象化的数据，是map[string]interface{}，这种很难用，这也是dynamic client原始特色，我觉得是没有schema导致的，导致不能对象化
	elbService, errCrd := elbServiceClient.Namespace("calmwu-namespace").Get("example-elbservice", metav1.GetOptions{})
	if errCrd != nil {
		calm_utils.Fatalf("Get ELBService CRD failed. err: %v", errCrd)
	}

	calm_utils.Debugf("elbService %s", litter.Sdump(elbService))

	// resourceVersion每次update都会变化，当你用dynamic client去update时要使用该rv
	calm_utils.Debugf("ELBService resource version: %s", elbService.GetResourceVersion())
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
