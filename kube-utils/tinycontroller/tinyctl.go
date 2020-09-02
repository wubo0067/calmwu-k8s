/*
 * @Author: calm.wu
 * @Date: 2020-09-02 15:03:34
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-02 18:01:21
 */

// Package tinycontroller  。。。
package tinycontroller

import (
	"context"
	"time"

	"kube-utils/vendor/k8s.io/apimachinery/pkg/util/wait"
	"kube-utils/vendor/k8s.io/client-go/util/workqueue"

	"github.com/sanity-io/litter"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	informappsv1 "k8s.io/client-go/informers/apps/v1"
	informcorev1 "k8s.io/client-go/informers/core/v1"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

type ResourceControllerOption func(*ResourceControllerOptions)

type ResourceControllerOptions struct {
	resourceType     ResourceType
	namespace        string
	kubeCfgPath      string
	resyncPeriod     time.Duration
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	threadiness      int
	// labelSelector string
	// fieldSelector string
}

type ResourceController struct {
	clientset   kubernetes.Interface
	queue       workqueue.RateLimitingInterface
	informer    cache.SharedIndexInformer
	threadiness int
}

var defaultResourceControllerOptions = ResourceControllerOptions{
	namespace: corev1.NamespaceAll,
}

// RunK8SResourceControllers 运行多个controller
func RunK8SResourceControllers(ctx context.Context, opts ...ResourceControllerOption) error {
	tcoptions := defaultResourceControllerOptions
	for _, o := range opts {
		o(&tcoptions)
	}

	klog.Infof("options: %s", litter.Sdump(tcoptions))

	// 获取k8sclientset
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = GetClientOutOfCluster(tcoptions.kubeCfgPath)
	} else {
		kubeClient = GetClient()
	}

	//stopCh := make(chan struct{})
	var informer cache.SharedIndexInformer

	// 根据类型构造infomer对象，调用client-go接口
	switch tcoptions.resourceType {
	case Pod:
		informer = informcorev1.NewFilteredPodInformer(kubeClient, tcoptions.namespace, tcoptions.resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, tcoptions.tweakListOptions)
	case Deployment:
		informer = informappsv1.NewFilteredDeploymentInformer(kubeClient, tcoptions.namespace, tcoptions.resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, tcoptions.tweakListOptions)
	default:
		return ErrResourceNotSupport
	}

	// 构造一个controller
	c := newResourceController(kubeClient, informer, &tcoptions)
	stopCh := make(chan struct{})
	defer close(stopCh)
	// 运行
	go c.Run(stopCh)

	// 等待停止
	<-ctx.Done()
	return nil
}

func newResourceController(client kubernetes.Interface, informer cache.SharedIndexInformer, tcoptions *ResourceControllerOptions) *ResourceController {
	return &ResourceController{
		threadiness: tcoptions.threadiness,
	}
}

func (rc *ResourceController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer rc.queue.ShutDown()

	klog.Info("Starting resource controller")

	// 启动informer
	go rc.informer.Run(stopCh)

	// 等待同步，每个资源一个，不像传统的controller多个资源用这个来同步
	if !cache.WaitForCacheSync(stopCh, rc.HasSynced) {
		utilruntime.HandleError(ErrCacheSyncTimeout)
		return
	}

	klog.Info("resource controller synced and ready")

	for i := 0; i < rc.threadiness; i++ {
		go wait.Until(rc.runWorker, time.Second, stopCh)
	}
}

// HasSynced is required for the cache.Controller interface.
func (rc *ResourceController) HasSynced() bool {
	return rc.informer.HasSynced()
}

func (rc *ResourceController) runWorker() {
	for rc.processNextWorkItem() {
	}
}

func (rc *ResourceController) processNextWorkItem() bool {
	return true
}
