/*
 * @Author: calm.wu
 * @Date: 2020-09-07 11:37:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 11:52:02
 */

package tinycontroller

import (
	"time"

	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
)

type ResourceControllerOption func(*ResourceControllerOptions)

func ResType(resourceType ResourceType) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.resourceType = resourceType
	}
}

func NS(namespace string) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.namespace = namespace
	}
}

func ResyncPeriod(period time.Duration) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.resyncPeriod = period
	}
}

func Threadiness(threadiness int) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.threadiness = threadiness
	}
}

func KubeCfg(kubeCfg string) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.kubeCfgPath = kubeCfg
	}
}

func Processor(resourceProcessor ResourceProcessor) ResourceControllerOption {
	return func(rcOptions *ResourceControllerOptions) {
		rcOptions.resourceProcessor = resourceProcessor
	}
}

// ListOption 设置lableSelector和fieldSelector的选择函数
func ListOption(tweakListOptions internalinterfaces.TweakListOptionsFunc) ResourceControllerOption {
	return func(rcOption *ResourceControllerOptions) {
		rcOption.tweakListOptions = tweakListOptions
	}
}
