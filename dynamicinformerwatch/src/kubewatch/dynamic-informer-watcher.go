/*
 * @Author: CALM.WU
 * @Date: 2021-04-14 17:01:28
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-14 18:00:56
 */

// Package kubewatch watch k8s resource
package kubewatch

import (
	"dynifr-watchres/src/config"

	"github.com/sanity-io/litter"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// DynamicInformerWatchResources watch kubernetes resources
func DynamicInformerWatchResources(dc dynamic.Interface, stopCh <-chan struct{}) error {
	cfgData := config.GetConfData()
	watchNamespace := cfgData.Namespace
	if watchNamespace == "all" {
		watchNamespace = v1.NamespaceAll
	}

	// it will give us back an informer for watch resource
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, watchNamespace, nil)

	watchResources := cfgData.WatchResources
	for index, resource := range watchResources {
		gvr, gs := schema.ParseResourceArg(resource.ResGVK)
		calmUtils.Debugf("%d watch resource:%s gvr:%s gs:%s ", index, resource, litter.Sdump(gvr), litter.Sdump(gs))

		// 传入资源的gvr，类似envoyfilters.v1alpha3.networking.istio.io，应该支持crd
		// 通过grv创建资源的通用informer对象
		ifr := f.ForResource(*gvr)
		// 获得informer的cache.SharedIndexInformer接口对象
		go startWatchingResource(ifr.Informer(), stopCh)
	}

	return nil
}

func startWatchingResource(s cache.SharedIndexInformer, stopCh <-chan struct{}) {
	// 注册回调事件
	evtHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			switch v := obj.(type) {
			case *unstructured.Unstructured:
				switch v.GroupVersionKind() {
				// vendor\k8s.io\client-go\informers\generic.go
				// resources 是小写复数
				// kind是大写驼峰
				case corev1.SchemeGroupVersion.WithKind("ConfigMap"):
					cm := &corev1.ConfigMap{}
					runtime.DefaultUnstructuredConverter.FromUnstructured(v.UnstructuredContent(), cm)
					calmUtils.Debugf("<AddEvt>. configmap convert unstruct to object, %s", litter.Sdump(cm))
				default:
					calmUtils.Debugf("<AddEvt>. name: %s, namespace: %s, gvk: %s", v.GetName(), v.GetNamespace(), v.GroupVersionKind().String())
				}
			case *unstructured.UnstructuredList:
				calmUtils.Debugf("<AddEvt>. list gvk: %s", v.GroupVersionKind().String())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ou := oldObj.(*unstructured.Unstructured)
			nu := newObj.(*unstructured.Unstructured)

			calmUtils.Debugf("<UpdateEvt>. ou: %s\n ===>\n nu: %s", litter.Sdump(ou), litter.Sdump(nu))
		},
		DeleteFunc: func(obj interface{}) {
			switch v := obj.(type) {
			case *unstructured.Unstructured:
				calmUtils.Debugf("<DeleteEvt>. u: %s", litter.Sdump(v))
			case *unstructured.UnstructuredList:
				calmUtils.Debugf("<DeleteEvt>. list u: %s", litter.Sdump(v))
			}
		},
	}
	s.AddEventHandler(evtHandlers)
	// 这里就是启动list/watch，
	s.Run(stopCh)
}
