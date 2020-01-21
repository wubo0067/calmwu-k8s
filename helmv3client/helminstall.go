/*
 * @Author: calm.wu
 * @Date: 2019-12-26 19:52:39
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-21 14:23:33
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"

	"github.com/snwfdhmp/errlog"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	diskcached "k8s.io/client-go/discovery/cached/disk"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"

	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	kubeCfgFile = "/root/.kube/config_bak"
)

var (
	errJudgement = errlog.NewLogger(&errlog.Config{
		// PrintFunc is of type `func (format string, data ...interface{})`
		// so you can easily implement your own logger func.
		// In this example, logrus is used, but any other logger can be used.
		// Beware that you should add '\n' at the end of format string when printing.
		PrintFunc:          calm_utils.Debugf,
		PrintSource:        true,  //Print the failing source code
		LinesBefore:        2,     //Print 2 lines before failing line
		LinesAfter:         1,     //Print 1 line after failing line
		PrintError:         true,  //Print the error
		PrintStack:         false, //Don't print the stack trace
		ExitOnDebugSuccess: false, //Exit if err
	})
)

type HelmClientConfigGetter struct {
	kubeconfigGetter clientcmd.KubeconfigGetter
}

func (g *HelmClientConfigGetter) Load() (*clientcmdapi.Config, error) {
	return g.kubeconfigGetter()
}

func (g *HelmClientConfigGetter) GetLoadingPrecedence() []string {
	return nil
}
func (g *HelmClientConfigGetter) GetStartingConfig() (*clientcmdapi.Config, error) {
	return g.kubeconfigGetter()
}
func (g *HelmClientConfigGetter) GetDefaultFilename() string {
	return ""
}
func (g *HelmClientConfigGetter) IsExplicitFile() bool {
	return false
}
func (g *HelmClientConfigGetter) GetExplicitFile() string {
	return ""
}
func (g *HelmClientConfigGetter) IsDefaultConfig(config *restclient.Config) bool {
	return false
}

var _ clientcmd.ClientConfigLoader = &HelmClientConfigGetter{}

type HelmConfigFlags struct {
	KubeCfgContent []byte
	genericclioptions.ConfigFlags
}

func (f *HelmConfigFlags) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	// config, err := clientcmd.NewClientConfigFromBytes(f.KubeCfgContent)
	// if err != nil {
	// 	calm_utils.Fatalf("NewClientConfigFromBytes failed. err:%s", err.Error())
	// }
	// calm_utils.Debugf("ToRawKubeConfigLoader config: %#v", config)

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&HelmClientConfigGetter{
			kubeconfigGetter: func() (*clientcmdapi.Config, error) {
				return clientcmd.Load(f.KubeCfgContent)
			},
		},
		// &clientcmd.ConfigOverrides{ClusterDefaults: clientcmdapi.Cluster{Server: "https://192.168.2.128:6443"},
		// 	ClusterInfo: clientcmdapi.Cluster{Server: "https://192.168.2.128:6443"}},
		&clientcmd.ConfigOverrides{ClusterDefaults: clientcmdapi.Cluster{Server: ""}},
	)
	//rawConfig, _ := config.RawConfig()
	clientConfig, _ := config.ClientConfig()
	calm_utils.Debugf("-----ToRawKubeConfigLoader-----, ClientConfig.Host:%s", clientConfig.Host)
	return config
}

func (f *HelmConfigFlags) ToRESTConfig() (*restclient.Config, error) {
	return f.ToRawKubeConfigLoader().ClientConfig()
}

func (f *HelmConfigFlags) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom resources) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = 100

	// retrieve a user-provided value for the "cache-dir"
	// defaulting to ~/.kube/http-cache if no user-value is given.
	httpCacheDir := filepath.Join(homedir.HomeDir(), ".kube", "http-cache")
	// if f.CacheDir != nil {
	// 	httpCacheDir = *f.CacheDir
	// }

	schemelessHost := strings.Replace(strings.Replace(config.Host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	discoveryCacheDir := filepath.Join(filepath.Join(homedir.HomeDir(), ".kube", "cache", "discovery"), schemelessHost)
	// discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(homedir.HomeDir(), ".kube", "cache", "discovery"), config.Host)
	calm_utils.Debugf("-----ToDiscoveryClient-----, discoveryCacheDir:%s", discoveryCacheDir)
	return diskcached.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, time.Duration(10*time.Minute))
}

func (f *HelmConfigFlags) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

var _ genericclioptions.RESTClientGetter = &HelmConfigFlags{}

func debug(format string, v ...interface{}) {
	calm_utils.Debug(fmt.Sprintf(format, v...))
}

func checkDeploymentTemplateFile(fileContent []byte) bool {
	return false
}

func checkServiceTemplateFile(fileContent []byte) bool {
	return false
}

func patchSciAnnotation(fileName string, fileContent []byte) ([]byte, error) {
	lineKindTag := -1
	lineSpecTag := -1
	lineTemplateTag := -1
	lineMetadataTag := -1

	newTemplateBuf := new(bytes.Buffer)

	lineNum := 0
	scanner := bufio.NewScanner(bytes.NewBuffer(fileContent))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineContent := scanner.Bytes()
		newTemplateBuf.Write(lineContent)
		newTemplateBuf.WriteByte('\n')
		//calm_utils.Debugf("%d\t%s", lineNum, lineContent)
		lineContentStr := calm_utils.Bytes2String(lineContent)
		if strings.Compare(lineContentStr, "kind: Deployment") == 0 {
			//calm_utils.Debug("this is a deployment yaml file")
			lineKindTag = lineNum
		}

		if strings.Compare(lineContentStr, "spec:") == 0 {
			//calm_utils.Debug("---find spec yaml node")
			lineSpecTag = lineNum
		}

		if strings.Compare(lineContentStr, "  template:") == 0 {
			//calm_utils.Debug("---find template yaml node")
			lineTemplateTag = lineNum
		}

		if strings.Compare(lineContentStr, "    metadata:") == 0 {
			//calm_utils.Debug("---find metadata yaml node")
			lineMetadataTag = lineNum
		}

		if lineMetadataTag > lineTemplateTag && lineTemplateTag > lineSpecTag && lineSpecTag > lineKindTag {
			//calm_utils.Debug("--->deployment template now will add sci annotation<---")
			lineKindTag = -1
			lineSpecTag = -1
			lineTemplateTag = -1
			lineMetadataTag = -1
			newTemplateBuf.WriteString(additionalAnnotations)
		}
		lineNum++
	}

	return newTemplateBuf.Bytes(), nil
}

func replaceDefaultClientConfig() {
	kubeCfgContent, err := ioutil.ReadFile(kubeCfgFile)
	if err != nil {
		calm_utils.Fatalf("readfile:%s failed. err:%s", err.Error())
	}

	apiCfg, err := clientcmd.Load(kubeCfgContent)
	if err != nil {
		calm_utils.Fatalf("clientcmd.Load failed. err:%s", err.Error())
	}

	clientcmd.DefaultClientConfig = *(clientcmd.NewDefaultClientConfig(*apiCfg, &clientcmd.ConfigOverrides{
		ClusterDefaults: clientcmdapi.Cluster{Server: "https://192.168.2.128:6443"},
		ClusterInfo:     clientcmdapi.Cluster{Server: "https://192.168.2.128:6443"}}).(*clientcmd.DirectClientConfig))
}

func helmInstall() {
	calm_utils.Debug("\n----------------helmInstall----------------")

	//replaceDefaultClientConfig()
	kubeCfgContent, err := ioutil.ReadFile(kubeCfgFile)
	//if err != nil {
	if errJudgement.Debug(err) {
		calm_utils.Fatalf("readfile:%s failed. err:%s", kubeCfgFile, err.Error())
	}
	getter := &HelmConfigFlags{
		KubeCfgContent: kubeCfgContent,
		ConfigFlags:    *genericclioptions.NewConfigFlags(true),
	}
	//getter.ConfigFlags = *genericclioptions.NewConfigFlags(true)
	namespace := "myguest"
	getter.ConfigFlags.Namespace = &namespace

	actionConfig := new(action.Configuration)
	//getter := kube.GetConfig(kubeCfgFile, "", "default")
	kc := kube.New(getter)

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		calm_utils.Fatalf("KubernetesClientSet failed. err:%s", err.Error())
	}

	nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		calm_utils.Fatalf("node list failed. err:%s", err.Error())
	}

	for index := range nodeList.Items {
		calm_utils.Debugf("node[%d]: %v", index, nodeList.Items[index])
	}

	// 使用secret，这里要保证namespace一致
	var store *storage.Storage
	d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
	d.Log = debug
	store = storage.Init(d)

	actionConfig.Releases = store
	actionConfig.KubeClient = kc
	actionConfig.RESTClientGetter = getter
	actionConfig.Log = debug

	installAction := action.NewInstall(actionConfig)
	// 在这里设置namespace
	installAction.Namespace = namespace
	installAction.ReleaseName = "myguestbook"

	sciVals := map[string]interface{}{
		"Network": map[string]interface{}{
			"RegionID": "a-b-c-d",
		},
	}

	// 先load本地的chart
	chart, err := loader.Load(guestBookChartDir)
	if err != nil {
		calm_utils.Fatalf("load chart:%s failed. err:%s", guestBookChartDir, err.Error())
	}

	//calm_utils.Debugf("chart.Templates:%v", chart.Templates)

	// 修改某一个chart文件内容
	for index := range chart.Templates {
		calm_utils.Debugf("template file name:%s", chart.Templates[index].Name)
		if chart.Templates[index].Name == "templates/guestbook-deployment.yaml" ||
			chart.Templates[index].Name == "templates/redis-master-deployment.yaml" {
			//
			chart.Templates[index].Data, _ = patchSciAnnotation(chart.Templates[index].Name, chart.Templates[index].Data)
			calm_utils.Debugf("patchSciAnnotation[%s] %s", chart.Templates[index].Name,
				calm_utils.Bytes2String(chart.Templates[index].Data))
		}
	}

	res, err := installAction.Run(chart, sciVals)
	if errJudgement.Debug(err) {
		calm_utils.Fatalf("run installAction failed. err:%s", err.Error())
	}

	calm_utils.Debugf("run installAction res:%v", res)
}
