/*
 * @Author: calmwu
 * @Date: 2020-09-13 19:49:03
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-09-13 20:38:56
 */

// https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go

package client

import (
	"testing"

	"golang.org/x/net/context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDynamicClient(t *testing.T) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		t.Fatalf("build restCfg failed. err:%s", err.Error())
	}

	// 通过config构造dynamic client
	dynamicClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		t.Fatalf("create dynamic client failed. err:%s", err.Error())
	}

	// gvr 这个是关键，可以通过gvk来获取gvr,
	// mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version) mapping.Resource
	// https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go

	// 通过schema.GroupVersionResource设置请求的资源版本和资源组，设置命名空间和请求参数,得到unstructured.UnstructuredList指针类型的PodLis
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unstructObj, err := dynamicClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{Limit: 500})
	if err != nil {
		t.Fatalf("use dynamic client list pods failed. err:%s", err.Error())
	}

	// 通过runtime.DefaultUnstructuredConverter函数将unstructured.UnstructuredList转为PodList类型
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		t.Fatalf("runtime FromUnstructured failed. err:%s", err.Error())
	}

	for index := range podList.Items {
		t.Logf("[%d] pod:%#v", index, podList.Items[index])
	}
}
