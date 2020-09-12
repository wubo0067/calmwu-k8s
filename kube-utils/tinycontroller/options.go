/*
 * @Author: calm.wu
 * @Date: 2020-09-07 11:37:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 11:52:02
 */

package tinycontroller

import "time"

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
