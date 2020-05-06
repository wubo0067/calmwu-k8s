/*
 * @Author: calm.wu
 * @Date: 2020-04-08 14:00:24
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-04-08 15:22:23
 */

package main

import (
	"strings"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/releaseutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
)

const notesFileSuffix = "NOTES.txt"

func renderTemplate(chartFullPath string, releaseName string, namespace string) {

	cfg, err := makeHelmConfiguration(namespace)
	if err != nil {
		calm_utils.Fatalf("make config failed. err:%s", err.Error())
	}

	dc, err := cfg.RESTClientGetter.ToDiscoveryClient()
	if err != nil {
		calm_utils.Fatalf("could not get Kubernetes discovery client, err:%s", err.Error())
	}

	dc.Invalidate()
	kubeVersion, err := dc.ServerVersion()
	if err != nil {
		calm_utils.Fatalf("could not get server version from Kubernetes. err:%s", err.Error())
	}

	apiVersions, err := action.GetVersionSet(dc)
	if err != nil {
		calm_utils.Fatalf("could not get apiVersions from Kubernetes. err:%s", err.Error())
	}

	caps := &chartutil.Capabilities{
		APIVersions: apiVersions,
		KubeVersion: chartutil.KubeVersion{
			Version: kubeVersion.GitVersion,
			Major:   kubeVersion.Major,
			Minor:   kubeVersion.Minor,
		},
	}
	calm_utils.Debugf("Capabilities:%s", litter.Sdump(caps))

	//----------------------------------------------------------------------------------

	sciVals := map[string]interface{}{
		"Network": map[string]interface{}{
			"RegionID": "a-b-c-d",
		},
	}

	// 加载本地的chart
	chart, err := loader.Load(chartFullPath)
	if err != nil {
		calm_utils.Fatalf("load chart:%s failed. err:%s", chartFullPath, err.Error())
	}

	// rest, err := cfg.RESTClientGetter.ToRESTConfig()
	// if err != nil {
	// 	calm_utils.Fatal("RESTClientGetter.ToRESTConfig failed. err:%s", err.Error())
	// }

	options := chartutil.ReleaseOptions{
		Name:      releaseName,
		Namespace: namespace,
		Revision:  1,
		IsInstall: true,
		IsUpgrade: false,
	}
	valuesToRender, err := chartutil.ToRenderValues(chart, sciVals, options, caps)
	if err != nil {
		calm_utils.Fatalf("chartutil ToRenderValues failed. err:%s", err.Error())
	}

	//files, err := engine.RenderWithClient(chart, valuesToRender, rest)
	files, err := engine.Render(chart, valuesToRender)
	if err != nil {
		calm_utils.Fatalf("engine.Render failed. err:%s", err.Error())
	}
	calm_utils.Debugf("files:%#v", files)

	// 先过滤掉NOTES.txt文件
	for k, _ := range files {
		calm_utils.Debugf("engine.Render k:[%s] ", k)
		if strings.HasSuffix(k, notesFileSuffix) {
			calm_utils.Debugf("filter out file:[%s]", k)
			delete(files, k)
			continue
		}
	}

	// 这个被拆分后的
	_, manifests, err := releaseutil.SortManifests(files, caps.APIVersions, releaseutil.InstallOrder)
	if err != nil {
		calm_utils.Fatalf("sort manifests failed. err:%s", err.Error())
	}

	for _, m := range manifests {
		checkManifest(m.Name, m.Content)
	}

	// // 查看清单
	// var manifests bytes.Buffer
	// // 写入数据
	// for _, f := range rel.Chart.CRDs() {
	// 	fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", f.Name, f.Data)
	// }
	// fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))

	// // 切分多个清单
	// splitManifests := releaseutil.SplitManifests(manifests.String())
	// for index, subManifest := range splitManifests {
	// 	//calm_utils.Debugf("[%s] subManifest----\n%s", index, subManifest)
	// 	checkManifest(index, subManifest)
	// }
}

func checkManifest(index string, subManifest string) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	runtimeObj, gvk, err := decode(calm_utils.String2Bytes(subManifest), nil, nil)

	if err != nil {
		calm_utils.Fatalf("[%s] scheme.Codecs.UniversalDeserializer().Decode failed, err:%s", index, err.Error())
	}

	switch obj := runtimeObj.(type) {
	case *v1.Service:
		calm_utils.Debugf("[%s] Service gvk[%s] name:%s Spec.type:%s", index, gvk.String(), obj.GetName(), obj.Spec.Type)
	case *appsv1.Deployment:
		calm_utils.Debugf("[%s] appsv1.Deployment gvk[%s] name:%s", index, gvk.String(), obj.GetName())
		if obj.Spec.Replicas != nil && *obj.Spec.Replicas > 70 {
			calm_utils.Warnf("current replicas:%d, the replicas limit is 70", *obj.Spec.Replicas)
		}
	case *extensionsv1beta1.Deployment:
		calm_utils.Debugf("[%s] extensionsv1beta1.Deployment gvk[%s] name:%s", index, gvk.String(), obj.GetName())
		if obj.Spec.Replicas != nil && *obj.Spec.Replicas > 70 {
			calm_utils.Warnf("current replicas:%d, the replicas limit is 70", *obj.Spec.Replicas)
		}
	case *extensionsv1beta1.Ingress:
		calm_utils.Warnf("[%s] extensionsv1beta1.Ingress gvk[%s] name:%s not support!", index, gvk.String(), obj.GetName())
	default:
		calm_utils.Warnf("[%s] Unknown gvk[%s] APIResource Type", index, gvk.String())
	}
}
