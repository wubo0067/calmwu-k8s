/*
 * @Author: calm.wu
 * @Date: 2020-09-07 11:52:31
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 14:28:33
 */

package tinycontroller

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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
		Processor(
			func(clientSet kubernetes.Interface, indexer cache.Indexer, key string, resourceType ResourceType) error {
				deploymentItfs, err := indexer.ByIndex(cache.NamespaceIndex, "test-indexer")

				if err != nil {
					err = errors.Wrap(err, "get deployments in test-indexer namespace failed.")
					klog.Error(err.Error())
					return err
				}

				if len(deploymentItfs) > 0 {
					for index := range deploymentItfs {
						if deployment, ok := deploymentItfs[index].(*appsv1.Deployment); ok {
							klog.Infof("index:%d deployment:%s", index,
								fmt.Sprintf("%s/%s", deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name))
						}
					}
				}

				return nil
			}),
	)
	if err != nil {
		klog.Info(err.Error())
	}
}
