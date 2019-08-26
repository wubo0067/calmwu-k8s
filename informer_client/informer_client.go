/*
 * @Author: calm.wu
 * @Date: 2019-08-14 17:02:36
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-14 17:58:09
 */

package main

import (
	"flag"
	"log"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func makeK8sClientSet() *kubernetes.Clientset {
	var kubeConfigFileName *string

	if home := homedir.HomeDir(); home != "" {
		kubeConfigFileName = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfigFileName = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	log.Printf("kubeConfigFileName:%s\n", *kubeConfigFileName)

	conf, err := clientcmd.BuildConfigFromFlags("", *kubeConfigFileName)
	if err != nil {
		log.Panic(err.Error())
	}

	// 构造clientset
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Panic(err.Error())
	}

	log.Println("make Clientset successed")
	return clientSet
}

func watchDeployment(clientSet *kubernetes.Clientset, stopCh chan struct{}) {
	// 生成一个新的factory
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	go factory.Start(stopCh)

	dpInformer := factory.Apps().V1().Deployments()
	informer := dpInformer.Informer()

	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Fatal("Timed out waiting for caches to sync")
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			dp := obj.(*appsv1.Deployment)
			log.Printf("new deployment:%s\n", dp.Name)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			newDP := newObj.(*appsv1.Deployment)
			oldDP := oldObj.(*appsv1.Deployment)
			log.Printf("update deployment:%s ===> %s \n", oldDP.Name, newDP.Name)
		},
		DeleteFunc: func(obj interface{}) {
			delDeployment := obj.(*appsv1.Deployment)
			log.Printf("del deployment:%s\n", delDeployment.Name)
		},
	})

	<-stopCh
}

func main() {
	stopInformerCh := make(chan struct{})
	defer close(stopInformerCh)

	clientSet := makeK8sClientSet()
	//
	go watchDeployment(clientSet, stopInformerCh)
	// 初始化informer
	factory := informers.NewSharedInformerFactory(clientSet, 0)

	// 启动informer， list & watch
	go factory.Start(stopInformerCh)

	podInformer := factory.Core().V1().Pods()
	informer := podInformer.Informer()

	// 从apiserver同步资源， list
	if !cache.WaitForCacheSync(stopInformerCh, informer.HasSynced) {
		log.Fatal("Timed out waiting for caches to sync")
	}

	// 使用自定义handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf("new pod:%s\n", pod.Name)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			oldPod := oldObj.(*corev1.Pod)
			log.Printf("update pod:%s ===> %s\n", oldPod.Name, newPod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			delPod := obj.(*corev1.Pod)
			log.Printf("del pod:%s\n", delPod.Name)
		},
	})

	// 创建lister
	podLister := podInformer.Lister()
	// 从podLister中获取pod
	podList, err := podLister.List(labels.Everything())
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("-------------------------pod list-----------------------------")
	for _, pod := range podList {
		log.Printf("----------Pod:%s status:%s\n", pod.Name, pod.Status.String())
	}
	<-stopInformerCh
	return
}
