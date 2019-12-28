/*
 * @Author: calm.wu
 * @Date: 2019-12-26 19:52:39
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-12-27 18:15:08
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	kubeCfgFile = "/root/.kube/config_bak"
)

type HelmConfigFlags struct {
	KubeCfgContent []byte
	genericclioptions.ConfigFlags
}

func (f *HelmConfigFlags) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	config, err := clientcmd.NewClientConfigFromBytes(f.KubeCfgContent)
	if err != nil {
		calm_utils.Fatalf("NewClientConfigFromBytes failed. err:%s", err.Error())
	}
	calm_utils.Debugf("ToRawKubeConfigLoader config: %#v", config)
	return config
}

func (f *HelmConfigFlags) ToRESTConfig() (*rest.Config, error) {
	return f.ToRawKubeConfigLoader().ClientConfig()
}

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
		calm_utils.Debugf("%d\t%s", lineNum, lineContent)
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
		ClusterInfo: clientcmdapi.Cluster{Server: ""},
	}).(*clientcmd.DirectClientConfig))
}

func helmInstall() {
	calm_utils.Debugf("\n----------------helmInstall----------------")

	//replaceDefaultClientConfig()
	kubeCfgContent, err := ioutil.ReadFile(kubeCfgFile)
	if err != nil {
		calm_utils.Fatalf("readfile:%s failed. err:%s", err.Error())
	}
	getter := &HelmConfigFlags{
		KubeCfgContent: kubeCfgContent,
	}
	getter.ConfigFlags = *genericclioptions.NewConfigFlags(true)
	namespace := "default"
	getter.Namespace = &namespace

	actionConfig := new(action.Configuration)
	//getter := kube.GetConfig(kubeCfgFile, "", "default")
	kc := kube.New(getter)

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		calm_utils.Fatalf("KubernetesClientSet failed. err:%s", err.Error())
	}

	_, err = clientset.CoreV1().Nodes().List(metav1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		calm_utils.Fatalf("node list failed. err:%s", err.Error())
	}

	// 使用configmap
	var store *storage.Storage
	d := driver.NewSecrets(clientset.CoreV1().Secrets("default"))
	d.Log = debug
	store = storage.Init(d)

	actionConfig.Releases = store
	actionConfig.KubeClient = kc
	actionConfig.RESTClientGetter = getter
	actionConfig.Log = debug

	installAction := action.NewInstall(actionConfig)
	installAction.Namespace = "default"
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
	if err != nil {
		calm_utils.Fatalf("run installAction failed. err:%s", err.Error())
	}

	calm_utils.Debugf("run installAction res:%v", res)
}
