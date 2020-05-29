/*
 * @Author: calm.wu
 * @Date: 2020-05-21 15:26:58
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-05-21 15:34:54
 */

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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	metricsHost               = "0.0.0.0"
	metricsPort         int32 = 8383
	operatorMetricsPort int32 = 8686
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

type SCIKube struct {
	config *rest.Config
	scheme *runtime.Scheme
	cache  cache.Cache
	client client.Client
}

func newSCIKube(config *rest.Config) *SCIKube {
	mapper, err := apiutil.NewDynamicRESTMapper(config)
	if err != nil {
		calm_utils.Fatalf("failed to call NewDynamicRESTMapper, err:%s", err.Error())
	}

	sciKubeScheme := scheme.Scheme

	cacheOptions := cache.Options{
		Scheme:    sciKubeScheme,
		Mapper:    mapper,
		Namespace: "",
	}

	// 创建cache
	sciKubeCache, err := cache.New(config, cacheOptions)
	if err != nil {
		calm_utils.Fatalf("failed to new Cache, err:%s", err.Error())
	}

	// 创建client
	c, err := client.New(config, client.Options{Scheme: sciKubeScheme, Mapper: mapper})
	if err != nil {
		calm_utils.Fatalf("failed to new Client, err:%s", err.Error())
	}

	sciKubeClient := client.DelegatingClient{
		Reader: &client.DelegatingReader{
			CacheReader:  sciKubeCache,
			ClientReader: c,
		},
		Writer:       c,
		StatusClient: c,
	}

	return &SCIKube{
		config: config,
		scheme: sciKubeScheme,
		cache:  sciKubeCache,
		client: sciKubeClient,
	}
}

func (sk *SCIKube) GetScheme() *runtime.Scheme {
	return sk.scheme
}

func (sk *SCIKube) Start(stopCh <-chan struct{}) error {
	go func() {
		if err := sk.cache.Start(stopCh); err != nil {
			calm_utils.Fatalf("failed to start Cache. err:%s", err.Error())
		}
	}()

	sk.cache.WaitForCacheSync(stopCh)
	calm_utils.Debug("wait for cache sync")
	return nil
}

func (sk *SCIKube) GetClient() client.Client {
	return sk.client
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

	sciKube := newSCIKube(kubeConfig)

	if err := elbservicescheme.AddToScheme(sciKube.GetScheme()); err != nil {
		calm_utils.Fatalf("failed to add elbservice scheme to runtime scheme: %s", err.Error())
	}

	stopCh := make(chan struct{})

	if err := sciKube.Start(stopCh); err != nil {
		calm_utils.Fatalf("sciKube start failed, err:%s", err.Error())
	}

	//time.Sleep(time.Second)

	/*
		(dlv) bt
		0  0x000000000170a9cb in sigs.k8s.io/controller-runtime/pkg/cache/internal.(*specificInformersMap).Get
			at /home/calmwu/Dev/k8s_space/elbservice-operator/vendor/sigs.k8s.io/controller-runtime/pkg/cache/internal/informers_map.go:168
		1  0x0000000001709e2a in sigs.k8s.io/controller-runtime/pkg/cache/internal.(*InformersMap).Get
			at /home/calmwu/Dev/k8s_space/elbservice-operator/vendor/sigs.k8s.io/controller-runtime/pkg/cache/internal/deleg_map.go:92
		2  0x0000000001710185 in sigs.k8s.io/controller-runtime/pkg/cache.(*informerCache).Get
			at /home/calmwu/Dev/k8s_space/elbservice-operator/vendor/sigs.k8s.io/controller-runtime/pkg/cache/informer_cache.go:60
		3  0x00000000016f4de9 in sigs.k8s.io/controller-runtime/pkg/client.(*DelegatingReader).Get
			at /home/calmwu/Dev/k8s_space/elbservice-operator/vendor/sigs.k8s.io/controller-runtime/pkg/client/split.go:51
		4  0x00000000016fc2eb in sigs.k8s.io/controller-runtime/pkg/client.(*DelegatingClient).Get
			at <autogenerated>:1
		5  0x00000000018c909a in main.main
			at ./access.go:176
		6  0x000000000043e748 in runtime.main
			at /usr/local/go/src/runtime/proc.go:203
		7  0x0000000000471961 in runtime.goexit
			at /usr/local/go/src/runtime/asm_amd64.s:1373

		//-------------------------------------------------------------------------------------------/
			229:	// newListWatch returns a new ListWatch object that can be used to create a SharedIndexInformer.
			230:	func createStructuredListWatch(gvk schema.GroupVersionKind, ip *specificInformersMap) (*cache.ListWatch, error) {
			231:		// Kubernetes APIs work against Resources, not GroupVersionKinds.  Map the
			232:		// groupVersionKind to the Resource API we will use.
			233:		mapping, err := ip.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		=> 234:		if err != nil {
			235:			return nil, err
			236:		}
			237:
			238:		client, err := apiutil.RESTClientForGVK(gvk, ip.config, ip.codecs)
			239:		if err != nil {
		(dlv) print mapping
		*k8s.io/apimachinery/pkg/api/meta.RESTMapping {
			Resource: k8s.io/apimachinery/pkg/runtime/schema.GroupVersionResource {
				Group: "apps",
				Version: "v1",
				Resource: "deployments",},
			GroupVersionKind: k8s.io/apimachinery/pkg/runtime/schema.GroupVersionKind {
				Group: "apps",
				Version: "v1",
				Kind: "Deployment",},
			Scope: k8s.io/apimachinery/pkg/api/meta.RESTScope(*k8s.io/apimachinery/pkg/api/meta.restScope) *{
				name: "namespace",},}

		//-------------------------------------------------------------------------------------------/
		(dlv) n
		> sigs.k8s.io/controller-runtime/pkg/cache/internal.createStructuredListWatch() /home/calmwu/Dev/k8s_space/elbservice-operator/vendor/sigs.k8s.io/controller-runtime/pkg/cache/internal/informers_map.go:249 (PC: 0x170bbdd)
		   244:		if err != nil {
		   245:			return nil, err
		   246:		}
		   247:
		   248:		// Create a new ListWatch for the obj
		=> 249:		return &cache.ListWatch{
		   250:			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		   251:				res := listObj.DeepCopyObject()
		   252:				isNamespaceScoped := ip.namespace != "" && mapping.Scope.Name() != meta.RESTScopeNameRoot
		   253:				err := client.Get().NamespaceIfScoped(ip.namespace, isNamespaceScoped).Resource(mapping.Resource.Resource).VersionedParams(&opts, ip.paramCodec).Do().Into(res)
		   254:				return res, err

	*/
	// 直接用Get测试下
	deployment := &appsv1.Deployment{}
	err = sciKube.GetClient().Get(context.TODO(), types.NamespacedName{Name: "example-appservice", Namespace: "default"}, deployment)
	if err != nil {
		calm_utils.Fatalf("failed to get deployment:[example-appservice] in namespace:[default], err: %s", err.Error())
	}
	calm_utils.Debugf("successed to deployment:[example-appservice] in namespace:[default] %s", litter.Sdump(deployment))

	deployment.ObjectMeta.Labels = map[string]string{"test": "update-by-controllerclient"}
	err = sciKube.GetClient().Update(context.TODO(), deployment)
	if err != nil {
		calm_utils.Fatalf("failed to update deployment:[example-appservice] in namespace:[default], err: %s", err.Error())
	}

	// 测试获取、修改cr
	exampleELBService := &elbserviceoperator.ELBService{}
	err = sciKube.GetClient().Get(context.TODO(), types.NamespacedName{Name: "example-elbservice", Namespace: "calmwu-namespace"}, exampleELBService)
	if err != nil {
		calm_utils.Fatalf("failed to elbservice:[example-elbservice] in namespace:[calmwu-namespace], err: %s", err.Error())
	}
	calm_utils.Debugf("successed to elbservice:[example-elbservice] in namespace:[calmwu-namespace] %s", litter.Sdump(exampleELBService))

	exampleELBService.Spec.Listener.Port = 4406
	err = sciKube.GetClient().Update(context.TODO(), exampleELBService)
	if err != nil {
		calm_utils.Fatalf("failed to update elbservice.Spec.Listener.Port, err:%s", err.Error())
	} else {
		calm_utils.Debugf("successed to update elbservice.Spec.Listener.Port")
	}

	elbServiceList := &elbserviceoperator.ELBServiceList{}
	err = sciKube.GetClient().List(context.TODO(), elbServiceList, client.InNamespace("calmwu-namespace"))
	if err != nil {
		calm_utils.Fatalf("failed to list elbservice, err:%s", err.Error())
	} else {
		calm_utils.Debugf("successed to list elbservice")
	}

	calm_utils.Debugf("elbServiceList:%s", litter.Sdump(elbServiceList))

	close(stopCh)
	time.Sleep(time.Second)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
