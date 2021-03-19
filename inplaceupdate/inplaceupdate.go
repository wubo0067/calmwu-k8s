/*
 * @Author: CALM.WU
 * @Date: 2021-03-19 16:31:02
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-03-19 17:36:11
 */

package main

import (
	"github.com/urfave/cli/v2"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	_clientSet *kubernetes.Clientset
)

func inplaceUpdatePod(ns, pod, newImage string) {
	calmUtils.Debugf("ns: %s, pod: %s, newImage: %s", ns, pod, newImage)
}

func main() {
	app := cli.NewApp()
	app.Name = "inplaceupdate-pod"
	app.Usage = "./inplaceupdate-pod --ns=xx --pod=yy --new-image=xxx.x.x"
	app.Action = func(c *cli.Context) error {
		inplaceUpdatePod(c.String("ns"), c.String("pod"), c.String("new-image"))
		return nil
	}
	app.Before = func(c *cli.Context) error {
		// 初始化clientset
		kubeCfg := c.String("kubeconfig")
		calmUtils.Debugf("kubeCfg: %s", kubeCfg)

		config, err := clientcmd.BuildConfigFromFlags("", kubeCfg)
		if err != nil {
			calmUtils.Error(err.Error())
			return err
		}

		_clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			calmUtils.Error(err.Error())
			return err
		}
		return nil
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "kubeconfig",
			Value: "/root/.kube/config",
			Usage: "k8s config",
		},
		&cli.StringFlag{
			Name:  "ns, n",
			Value: "default",
			Usage: "Set namespace name",
		},
		&cli.StringFlag{
			Name:  "pod, p",
			Value: "",
			Usage: "Set pod name",
		},
		&cli.StringFlag{
			Name:  "new-image, i",
			Value: "",
			Usage: "Set the image of the upgrade update",
		},
	}
}
