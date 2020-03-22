/*
 * @Author: calmwu
 * @Date: 2020-03-22 13:11:13
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-03-22 13:34:11
 */

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/micro/cli"
	uuid "github.com/satori/go.uuid"
	calmwu_utils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var (
	logger *log.Logger
)

func init() {
	logger = calmwu_utils.NewSimpleLog(nil)
}

func leaderRun(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			logger.Printf("Leader ticker working at :%v", t)
		case <-ctx.Done():
			logger.Printf("Leader run stop.")
		}
	}
}

func leaderElection(clientSet *kubernetes.Clientset) {
	hostName, _ := os.Hostname()
	id := uuid.NewV4()
	ownerID := hostName + "_" + id.String()

	logger.Printf("my id is %s", ownerID)

	rl, err := resourcelock.New("endpoints", // support endpoints and configmaps
		"default",
		"test-leaderelection",
		clientSet.CoreV1(),
		clientSet.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: ownerID,
		})
	if err != nil {
		logger.Fatalf("create ResourceLock error: %v", err)
	}

	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				logger.Println("you are the leader")
				leaderRun(ctx)
			},
			OnStoppedLeading: func() {
				logger.Fatalf("leaderelection lost")
			},
			OnNewLeader: func(identity string) {
				logger.Printf("New leader is %s", identity)
			},
		},
		Name: "test-leaderelection",
	})
}

func main() {
	app := cli.NewApp()
	app.Name = "leaderelection"
	app.Usage = "leader election"
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

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Fatal(err.Error())
		}

		leaderElection(clientSet)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Watch Deployment exit!")
}
