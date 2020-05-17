/*
 * @Author: calmwu
 * @Date: 2020-05-17 10:39:16
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-05-17 10:51:46
 */

// 使用dynclient "sigs.k8s.io/controller-runtime/pkg/client"去访问CRD资源
// 使用sigs.k8s.io/controller-runtime/pkg/client能更好的去CRUD CR资源，很好的结合GVK、GVR、RESTMapper、Scheme、CR反序列化，CRUD扩展k8s资源更方便

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	elbservicescheme "calmwu.org/elbservice-operator/pkg/apis"
	elbserviceoperator "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"

	extscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	cached "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/kubernetes"
	cgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// GetKubeconfigAndNamespace returns the *rest.Config and default namespace defined in the
// kubeconfig at the specified path. If no path is provided, returns the default *rest.Config
// and namespace
func GetKubeconfigAndNamespace(configPath string) (*rest.Config, string, error) {
	var clientConfig clientcmd.ClientConfig
	var apiConfig *clientcmdapi.Config
	var err error
	if configPath != "" {
		apiConfig, err = clientcmd.LoadFromFile(configPath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load user provided kubeconfig: %v", err)
		}
	} else {
		apiConfig, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
		if err != nil {
			return nil, "", fmt.Errorf("failed to get kubeconfig: %v", err)
		}
	}
	clientConfig = clientcmd.NewDefaultClientConfig(*apiConfig, &clientcmd.ConfigOverrides{})
	kubeconfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, "", err
	}
	return kubeconfig, namespace, nil
}

func main() {
	calm_utils.Debug("Starting test")

	var kubeConfigPath *string
	if home := homeDir(); home != "" {
		kubeConfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	kubeConfig, namespace, err := GetKubeconfigAndNamespace(*kubeConfigPath)
	if err != nil {
		calm_utils.Panic(err.Error())
	}

	calm_utils.Debugf("namespace:%s", namespace)

	// 根据kubeconfig构造kubeclient,
	kubeclient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		calm_utils.Fatalf("failed to build the kubeclient: %w", err)
	}

	// 将crd的schema加入构造的schema中，其实加入cgoscheme.Scheme就足够满足crd+基本类型了，这里演示下如果创建一个完整的scheme
	//scheme := cgoscheme.Scheme
	scheme := runtime.NewScheme()
	if err := cgoscheme.AddToScheme(scheme); err != nil {
		calm_utils.Fatalf("failed to add cgo scheme to runtime scheme: %w", err)
	}

	if err := extscheme.AddToScheme(scheme); err != nil {
		calm_utils.Fatalf("failed to add ext scheme to runtime scheme: %w", err)
	}

	if err := elbservicescheme.AddToScheme(scheme); err != nil {
		calm_utils.Fatalf("failed to add elbservice scheme to runtime scheme: %w", err)
	}

	// 为什么这里要构造一个discoveryClient，还是基于cache的
	cachedDiscoveryClient := cached.NewMemCacheClient(kubeclient.Discovery())

	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	restMapper.Reset()

	// 构造dynamic Client
	dynClient, err := dynclient.New(kubeConfig, dynclient.Options{Scheme: scheme, Mapper: restMapper})
	if err != nil {
		calm_utils.Fatalf("failed to build the dynamic client: %w", err)
	}

	// 立即测试这个dynClient
	err = wait.PollImmediate(time.Second, time.Second*10, func() (done bool, err error) {
		// 测试dynClient的list功能
		err = dynClient.List(context.TODO(), &elbserviceoperator.ELBServiceList{}, dynclient.InNamespace("calmwu-namespace"))
		if err != nil {
			calm_utils.Debug("failed to list ELBServiceList by dynamic client")
			restMapper.Reset()
			return false, nil
		}
		calm_utils.Debug("successed to list ELBServiceList by dynamic client")

		return true, nil
	})

	// 通过dynClient获取cr对象
	exampleELBService := &elbserviceoperator.ELBService{}
	err = dynClient.Get(context.TODO(), types.NamespacedName{Name: "example-elbservice", Namespace: "calmwu-namespace"}, exampleELBService)
	if err != nil {
		calm_utils.Fatalf("failed to get ELBService cr by dynamic client: %w", err)
	}

	calm_utils.Debugf("successed to get ELBService cr by dynamic client: %s", litter.Sdump(exampleELBService))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
