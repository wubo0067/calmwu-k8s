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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

// DynamicInformerWatchResources watch kubernetes resources
func DynamicInformerWatchResources(dc dynamic.Interface, stopCh <-chan struct{}) error {
	cfgData := config.GetConfData()
	watchNamespace := cfgData.Namespace
	if watchNamespace == "all" {
		watchNamespace = v1.NamespaceAll
	}

	// it will give us back an informer for watch resource
	_ = dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, watchNamespace, nil)

	watchResources := cfgData.WatchResources
	for index, resource := range watchResources {
		gvr, gs := schema.ParseResourceArg(resource.ResGVK)
		calmUtils.Debugf("%d watch resource:%s gvr:%s gs:%s ", index, resource, litter.Sdump(gvr), litter.Sdump(gs))
	}

	return nil
}
