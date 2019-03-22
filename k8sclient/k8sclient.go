/*
 * @Author: calm.wu
 * @Date: 2019-03-18 18:31:10
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-03-20 16:22:08
 */

package main

import (
	// "k8s.io/kubernetes/pkg/controller/deployment"
	// "k8s.io/kubernetes/pkg/registry/apps/deployment"
	// "k8s.io/kubernetes/pkg/kubectl/util/deployment"
	// "k8s.io/kubernetes/pkg/api/v1/service"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func listPod(clientSet *kubernetes.Clientset) {
	pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		logger.Fatal(err)
	}
	for i, pod := range pods.Items {
		if pod.Status.Phase == apiv1.PodPending {
			color.New(color.FgBlue).Printf("\t{%d}: pod:%s status:%s PodIP:%s\n", i, pod.Name, pod.Status.Phase, pod.Status.PodIP)
		} else {
			logger.Printf("\t{%d}: pod:%s status:%s PodIP:%s\n", i, pod.Name, pod.Status.Phase, pod.Status.PodIP)
		}
	}
}

func listDeployment(clientSet *kubernetes.Clientset) {
	// 由于我的yam文件中写的是apiVersion: extensions/v1beta1，所以这里使用ExtensionsV1beta1来查询，这点很重要
	deploymentsClient := clientSet.ExtensionsV1beta1().Deployments(apiv1.NamespaceDefault)
	//deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployments, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		logger.Fatal(err)
	}
	
	for i, deployment := range deployments.Items {
		logger.Printf("\t{%d}: deployment:%s (%d replicas)\n", i, deployment.Name, *deployment.Spec.Replicas)
	}
}

func listServices(clientSet *kubernetes.Clientset) {
	var kubeClient kubernetes.Interface = clientSet
	servicesClient := kubeClient.CoreV1().Services(apiv1.NamespaceDefault)

	services, err := servicesClient.List(metav1.ListOptions{})
	if err != nil {
		logger.Fatal(err)
	}
	
	for i := range services.Items {
		service := &services.Items[i]
		logger.Printf("\t{%d}: service:%s labes:%#v\n", i, service.Name, service.Labels)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "k8sclient"
	app.Usage = "k8sclient"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "kubeconfig",
			Value: "",
			Usage: "kubeconfig",
		},
	}

	app.Action = func(c *cli.Context) error {
		kubeconfig := c.String("kubeconfig")
		// 判断文件是否存在
		_, err := os.Stat(kubeconfig)
		if err != nil {
			logger.Fatal(err)
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Fatal(err)
		}

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Fatal(err)
		}

		listPod(clientSet)
		listDeployment(clientSet)
		listServices(clientSet)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("k8sclient exit!")
}
