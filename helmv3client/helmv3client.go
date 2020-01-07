/*
 * @Author: calm.wu
 * @Date: 2019-12-24 10:42:23
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-04 23:00:26
 */

package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

const (
	guestBookChartDir = "/home/calmwu/Dev/Downloads/workshop/kubecon2019china/charts/guestbook"
	chartFilePath     = "/home/calmwu/Dev/Downloads/workshop/kubecon2019china/charts/guestbook/Chart.yaml"
	valuesFilePath    = "/home/calmwu/Dev/Downloads/workshop/kubecon2019china/charts/guestbook/values.yaml"
	templateFilePath  = "/home/calmwu/Dev/Downloads/workshop/kubecon2019china/charts/guestbook/templates/guestbook-deployment.yaml"
	templateFilePath1 = "/home/calmwu/Dev/Downloads/workshop/kubecon2019china/charts/guestbook/templates/redis-master-deployment.yaml"
	yamlFilePath      = "/home/calmwu/Dev/Downloads/nginx-dp.yaml"
)

var additionalAnnotations = `      annotations:
        io.kubernetes-network.region-id: {{ .Values.Network.RegionID }}
        io.kubernetes.cri.untrusted-workload: "true"
`

func loadChartFromDir(chartDir string) {
	l, err := loader.Loader(chartDir)
	if err != nil {
		calm_utils.Fatalf("load chart from dir:%s failed. err:%s", chartDir, err.Error())
	}

	chart, err := l.Load()
	if err != nil {
		calm_utils.Fatalf("chart load failed. err:%s", err.Error())
	}

	for _, rawFile := range chart.Raw {
		calm_utils.Debugf("rawFile info:%+v", rawFile.Name)
	}

	calm_utils.Debugf("-------------------")

	for _, tempFile := range chart.Templates {
		calm_utils.Debugf("templateFile info:%+v", tempFile.Name)
	}

	calm_utils.Debugf("-------------------")

	for key, val := range chart.Values {
		calm_utils.Debugf("Values key:%s val:%v", key, val)
	}

	calm_utils.Debugf("-------------------")

	calm_utils.Debugf("chart MetaData:%#v", chart.Metadata)

	calm_utils.Debugf("chart %s load successed!", chartDir)
}

func loadChartFile(chartFilePath string) {
	calm_utils.Debugf("\n----------------loadChartFile----------------")
	chartMetadata, err := chartutil.LoadChartfile(chartFilePath)
	if err != nil {
		calm_utils.Fatalf("load Chart.yaml failed. err:%s", err.Error())
	}

	calm_utils.Debugf("chartMetadata %#v", chartMetadata)
}

func loadValuesFile(valuesFilePath string) {
	calm_utils.Debugf("\n----------------loadVauesFile----------------")
	values, err := chartutil.ReadValuesFile(valuesFilePath)
	if err != nil {
		calm_utils.Fatalf("load values.yaml failed. err:%s", err.Error())
	}

	for vKey, vVal := range values {
		calm_utils.Debugf("vKey[%s] vVal[%v] vValType:%s", vKey, vVal, reflect.TypeOf(vVal).String())
	}
}

func loadYamlFile(yamlFilePath string) {
	calm_utils.Debugf("\n----------------loadYamlFile----------------")
	templateData, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		calm_utils.Fatalf("load nginx-dp.yaml failed. err:%s", err.Error())
	}

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(templateData, &m); err != nil {
		calm_utils.Fatalf("yaml unmarshal nginx-dp.yaml failed. err:%s", err.Error())
	}

	calm_utils.Debugf("nginx-dp.yaml ===> %s", litter.Sdump(m))
}

func externTemplateFile(templateFilePath string) {
	calm_utils.Debugf("\n----------------loadTemplateFile----------------")
	templateData, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		calm_utils.Fatalf("load guestbook-deployment.yaml failed. err:%s", err.Error())
	}

	lineKindTag := -1
	lineSpecTag := -1
	lineTemplateTag := -1
	lineMetadataTag := -1

	newTemplateBuf := new(bytes.Buffer)

	lineNum := 0
	scanner := bufio.NewScanner(bytes.NewBuffer(templateData))
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
			calm_utils.Debug("--->deployment template now will add sci annotation<---")
			lineKindTag = -1
			lineSpecTag = -1
			lineTemplateTag = -1
			lineMetadataTag = -1
			newTemplateBuf.WriteString(additionalAnnotations)
		}
		lineNum++
	}

	calm_utils.Debugf("%s", newTemplateBuf.String())
}

func main() {
	calm_utils.Debug("helmV3 client start")
	//loadChartFromDir(guestBookChartDir)

	//loadYamlFile(yamlFilePath)

	//loadChartFile(chartFilePath)
	//loadValuesFile(valuesFilePath)
	//externTemplateFile(templateFilePath)
	//externTemplateFile(templateFilePath1)
	//helmInstall()
	//patchTemplateFileWithSCISpecific("./test_patch_template/calico_etcd.yaml")
	patchTemplateFileWithSCISpecific("./test_patch_template/metrics-deployment.yaml")
	patchTemplateFileWithSCISpecific("./test_patch_template/cerebro-deployment.yaml")
}
