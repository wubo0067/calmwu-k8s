/*
 * @Author: calm.wu
 * @Date: 2020-09-02 15:31:22
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-02 16:57:03
 */

package tinycontroller

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig(kubeCfgPath string) (*rest.Config, error) {
	kubeconfigPath := kubeCfgPath
	if kubeconfigPath == "" {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster(kubeCfgPath string) kubernetes.Interface {
	config, err := buildOutOfClusterConfig(kubeCfgPath)
	if err != nil {
		klog.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Can not get kubernetes config: %v", err)
	}

	return clientset
}
