/*
 * @Author: calmwu
 * @Date: 2020-09-13 19:24:14
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-09-13 20:59:39
 */

package client

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRestClient(t *testing.T) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		t.Fatalf("build restCfg failed. err:%s", err.Error())
	}

	// 配置API路径和请求的资源组/资源版本信息
	restCfg.APIPath = "api"
	restCfg.GroupVersion = &corev1.SchemeGroupVersion
	restCfg.NegotiatedSerializer = scheme.Codecs

	// RESTClient是最基础的客户端，对HTTP Request进行了封装，实现了RESTful风格的API，其他的三个client都是基于RESTClient实现的。
	restClient, err := rest.RESTClientFor(restCfg)
	if err != nil {
		t.Fatalf("create rest client failed. err:%s", err.Error())
	}

	podLst := &corev1.PodList{}
	err = restClient.Get().
		Namespace("").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(podLst)

	if err != nil {
		t.Fatalf("get pod list failed. err:%s", err.Error())
	}

	for index := range podLst.Items {
		t.Logf("[%d] pod:%#v", index, podLst.Items[index])
	}
}
