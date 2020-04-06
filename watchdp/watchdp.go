/*
 * @Author: calm.wu
 * @Date: 2020-02-20 15:34:17
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-02-20 15:44:14
 */

// Package main for watch Deployment
package main

import (
	"log"
	"os"
	"reflect"

	"github.com/sanity-io/litter"
	"github.com/urfave/cli"
	calmwu_utils "github.com/wubo0067/calmwu-go/utils"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	logger *log.Logger
)

func init() {
	logger = calmwu_utils.NewSimpleLog(nil)
}

func watchDeployment(clientSet *kubernetes.Clientset) {
	dpWatcher, err := clientSet.AppsV1().Deployments("default").Watch(metav1.ListOptions{
		LabelSelector: "app=test-scale-status",
	})

	if err != nil {
		logger.Fatal(err.Error())
	}

	for e := range dpWatcher.ResultChan() {
		logger.Printf("event type[%s] event object[%s:%T]", e.Type, reflect.TypeOf(e).String(), e.Object)

		if dpObj, ok := e.Object.(*appsv1.Deployment); ok {
			logger.Printf("Deployment ObjectMeta[%s] status:%s", litter.Sdump(dpObj.ObjectMeta),
				litter.Sdump(dpObj.Status))
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "watchdp"
	app.Usage = "Watch Deployment"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "kubeconfig",
			Value: "/root/.kube/config",
			Usage: "k8s config",
		},
	}

	app.Action = func(c *cli.Context) error {
		// 初始化log
		kubeconfig := c.String("kubeconfig")
		logger.Printf("kubeconfig:[%s]", kubeconfig)

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Printf("config:%s", litter.Sdump(config))
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Fatal(err.Error())
		}

		watchDeployment(clientSet)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Watch Deployment exit!")
}
