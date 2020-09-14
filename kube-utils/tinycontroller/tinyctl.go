/*
 * @Author: calm.wu
 * @Date: 2020-09-02 15:03:34
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 14:43:21
 */
// https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html

// Package tinycontroller  。。。
package tinycontroller

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informappsv1 "k8s.io/client-go/informers/apps/v1"
	informcorev1 "k8s.io/client-go/informers/core/v1"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type ResourceProcessor func(clientSet kubernetes.Interface, indexer cache.Indexer, key string, resourceType ResourceType) error

type ResourceControllerOptions struct {
	resourceType      ResourceType
	namespace         string
	kubeCfgPath       string
	resyncPeriod      time.Duration                           // ResyncPeriod defines how often the controller goes through all items remaining in the cache and fires the UpdateFunc again.
	tweakListOptions  internalinterfaces.TweakListOptionsFunc // 对象过滤，apiserver将过滤后的数据发送给cache
	threadiness       int                                     // 处理资源对象的routine数量
	resourceIndexName string                                  // 索引名字
	resourceIndexFunc cache.IndexFunc                         // 自定义索引函数，cache会将该函数作用于对象，返回对象的值，这个值决定了相同值的一堆对象
	resourceProcessor ResourceProcessor                       // 资源处理函数
}

type ResourceController struct {
	clientset         kubernetes.Interface
	queue             workqueue.RateLimitingInterface
	informer          cache.SharedIndexInformer
	threadiness       int
	resourceType      ResourceType
	resourceProcessor ResourceProcessor // 资源处理函数
}

const (
	maxRetries = 5
)

var defaultResourceControllerOptions = ResourceControllerOptions{
	namespace:    corev1.NamespaceAll,
	resyncPeriod: 0,
	threadiness:  1,
}

// RunK8SResourceControllers 运行多个controller
func RunK8SResourceControllers(stopCh <-chan struct{}, opts ...ResourceControllerOption) error {
	tcoptions := &defaultResourceControllerOptions
	for _, o := range opts {
		o(tcoptions)
	}

	klog.Infof("options: %#v", tcoptions)

	// 获取k8sclientset
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = GetClientOutOfCluster(tcoptions.kubeCfgPath)
	} else {
		kubeClient = GetClient()
	}

	// 每种资源都有自己的informer，crd资源也是如此
	var informer cache.SharedIndexInformer

	// 加速对象在cache中查询，rc.informer.GetIndexer
	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	if tcoptions.resourceIndexName != "" && tcoptions.resourceIndexFunc != nil {
		indexers[tcoptions.resourceIndexName] = tcoptions.resourceIndexFunc
	}

	// 根据类型构造infomer对象，调用client-go接口
	switch tcoptions.resourceType {
	case Pod:
		informer = informcorev1.NewFilteredPodInformer(kubeClient, tcoptions.namespace, tcoptions.resyncPeriod,
			indexers, tcoptions.tweakListOptions)
	case Deployment:
		informer = informappsv1.NewFilteredDeploymentInformer(kubeClient, tcoptions.namespace, tcoptions.resyncPeriod,
			indexers, tcoptions.tweakListOptions)
	default:
		return ErrResourceNotSupport
	}

	// 构造一个controller
	c, err := newResourceController(tcoptions.resourceType, kubeClient, informer, tcoptions)
	if err != nil {
		return err
	}

	// 运行
	c.Run(stopCh)
	return nil
}

func newResourceController(resourceType ResourceType, client kubernetes.Interface, informer cache.SharedIndexInformer, tcoptions *ResourceControllerOptions) (*ResourceController, error) {
	rc := &ResourceController{
		threadiness:       tcoptions.threadiness,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()), // 构造队列，存放key
		informer:          informer,
		clientset:         client,
		resourceType:      resourceType,
		resourceProcessor: tcoptions.resourceProcessor,
	}

	// 注册事件处理函数
	rc.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(obj); err != nil {
				utilruntime.HandleError(errors.Wrap(err, "informer add event handler get obj key failed."))
				return
			} else {
				klog.Infof("resource[%s] controller AddFunc add key:%s to workqueue", resourceType, key)
				rc.queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if oldObj.(*appsv1.Deployment).ResourceVersion == newObj.(*appsv1.Deployment).ResourceVersion {
				return
			}

			if key, err := cache.MetaNamespaceKeyFunc(newObj); err != nil {
				utilruntime.HandleError(errors.Wrapf(err, "resource[%s] controller informer update event handler get newObj key failed.", resourceType))
				return
			} else {
				klog.Infof("resource[%s] controller UpdateFunc add key:%s to workqueue", resourceType, key)
				rc.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// 处理delete传入的object，有两种可能，也许是你期望的类型，有可能是DeletedFinalStateUnknown类型
			var object metav1.Object
			var ok bool

			if object, ok = obj.(metav1.Object); !ok {
				if tombStone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
					if _, ok := tombStone.Obj.(metav1.Object); ok {
						klog.Infof("resource[%s] controller Recovered deleted object '%s' from tombstone", resourceType, object.GetName())

						if key, err := cache.MetaNamespaceKeyFunc(obj); err != nil {
							utilruntime.HandleError(errors.Wrapf(err, "resource[%s] controller informer delete event handler get newObj key failed.", resourceType))
							return
						} else {
							klog.Infof("resource[%s] controller DeleteFunc add key:%s to workqueue", resourceType, key)
							rc.queue.Add(key)
						}
					} else {
						utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
					}
				} else {
					utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
				}
			} else {
				if key, err := cache.MetaNamespaceKeyFunc(obj); err != nil {
					utilruntime.HandleError(errors.Wrap(err, "informer delete event handler get newObj key failed."))
					return
				} else {
					klog.Infof("resource[%s] controller DeleteFunc add key:%s to workqueue", resourceType, key)
					rc.queue.Add(key)
				}
			}
		},
	})

	return rc, nil
}

func (rc *ResourceController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer rc.queue.ShutDown()

	klog.Infof("Starting resource[%s] controller", rc.resourceType)

	// 启动informer
	go rc.informer.Run(stopCh)

	// 等待同步，每个资源一个，不像传统的controller多个资源用这个来同步
	if !cache.WaitForCacheSync(stopCh, rc.HasSynced) {
		utilruntime.HandleError(ErrCacheSyncTimeout)
		return
	}

	klog.Infof("resource[%s] controller synced and ready", rc.resourceType)

	// 每秒钟去执行一次runWorker
	for i := 0; i < rc.threadiness; i++ {
		go wait.Until(rc.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")
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
	// 从队列中获取，这个队列的特性要记住
	key, quit := rc.queue.Get()

	if quit {
		klog.V(7).Infof("resource[%s] controller workqueue be shutdown", rc.resourceType)
		return false
	}

	// 根据处理结果，处理key
	defer rc.queue.Done(key)
	err := rc.processItem(key.(string))
	if err == nil {
		// 处理成功
		klog.Infof("resource[%s] controller successfully processItem '%s'", rc.resourceType, key)
		rc.queue.Forget(key)
	} else if rc.queue.NumRequeues(key) < maxRetries {
		// 重新入队次数小于最大重试次数
		klog.Errorf("resource[%s] controller error processing %s (will retry): %v", rc.resourceType, key, err)
		rc.queue.AddRateLimited(key)
	} else {
		// 超过重试次数，放弃处理
		klog.Errorf("resource[%s] controller error processing %s (giving up): %v", rc.resourceType, key, err)
		rc.queue.Forget(key)
	}

	return true
}

func (rc *ResourceController) processItem(key string) error {
	if rc.resourceProcessor != nil {
		return rc.resourceProcessor(rc.clientset, rc.informer.GetIndexer(), key, rc.resourceType)
	}

	obj, exists, err := rc.informer.GetIndexer().GetByKey(key)

	if err != nil {
		err = errors.Wrapf(err, "resource[%s] controller get object by key:%s failed.", rc.resourceType, key)
		klog.Error(err.Error())
		return err
	}

	if !exists {
		// 如果是不存在，就需要忽略这个错误，不需要重新入队列
		// 删除对象事件会执行到这里
		klog.Warningf("resource[%s] controller object %s not found.", rc.resourceType, key)
		return nil
	}

	klog.Infof("Get Resource[%s] \nmeta: %s", key, litter.Sdump(GetObjectMetaData(obj)))
	return nil
}
