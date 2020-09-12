/*
 * @Author: calm.wu
 * @Date: 2020-09-07 11:52:31
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 14:28:33
 */

package tinycontroller

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/klog"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func SetupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		klog.Info("catch shutdown signal")
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func RunDeploymentController() {
	stopCh := SetupSignalHandler()

	err := RunK8SResourceControllers(stopCh,
		ResType(Deployment),
		Threadiness(1),
		KubeCfg("/root/.kube/config"),
		ResyncPeriod(15*time.Second),
	)
	if err != nil {
		klog.Info(err.Error())
	}
}
