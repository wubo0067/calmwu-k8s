/*
 * @Author: calm.wu
 * @Date: 2020-01-21 14:20:31
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-04-08 14:37:25
 */

// yaml编解码对象的相关文档
// https://github.com/kubernetes/apimachinery/blob/master/pkg/util/yaml/decoder.go
// https://github.com/kubernetes/client-go/issues/193
// https://github.com/appscode/voyager/issues/964

package main

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/scheme"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeHelmConfiguration(namespace string) (*action.Configuration, error) {
	kubeCfgContent, err := ioutil.ReadFile(kubeCfgFile)
	//if err != nil {
	if errJudgement.Debug(err) {
		calm_utils.Fatalf("readfile:%s failed. err:%s", kubeCfgFile, err.Error())
	}
	getter := &HelmConfigFlags{
		KubeCfgContent: kubeCfgContent,
		ConfigFlags:    *genericclioptions.NewConfigFlags(true),
	}
	getter.ConfigFlags.Namespace = &namespace

	//getter := kube.GetConfig(kubeCfgFile, "", "default")
	kc := kube.New(getter)

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		calm_utils.Fatalf("KubernetesClientSet failed. err:%s", err.Error())
	}

	// 使用secret，这里要保证namespace一致
	var store *storage.Storage
	d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
	d.Log = debug
	store = storage.Init(d)

	actionConfig := new(action.Configuration)
	actionConfig.Releases = store
	actionConfig.KubeClient = kc
	actionConfig.RESTClientGetter = getter
	actionConfig.Log = debug

	return actionConfig, nil
}

func status(releaseName string, namespace string) {
	helmActionConfig, _ := makeHelmConfiguration(namespace)

	statusAction := action.NewStatus(helmActionConfig)
	release, err := statusAction.Run(releaseName)
	if err != nil {
		calm_utils.Fatalf("run Action status failed, err:%s", err.Error())
	}

	calm_utils.Debug(release.Manifest)

	// from string to io.ReadCloser
	r := ioutil.NopCloser(strings.NewReader(release.Manifest))
	decoder := yaml.NewDocumentDecoder(r)
	defer decoder.Close()

	// 循环读取
	dataChunk := make([]byte, 128)
	yamlContent := []byte(nil)
	yamlContentLen := 0
	for {
		len, err := decoder.Read(dataChunk)
		//calm_utils.Debugf("read yaml data len:%d err:%v", len, err)
		if err == io.EOF {
			break
		}
		yamlContent = append(yamlContent, dataChunk...)
		yamlContentLen += len
		if len <= 128 && err == nil {
			// 读取完毕，ok这样就读取一个完整的yaml了
			//calm_utils.Debugf("read yaml total len:%d", yamlContentLen)
			//calm_utils.Debug(calm_utils.Bytes2String(yamlContent[:yamlContentLen]))
			decodeFromYamlContent(helmActionConfig.KubeClient.(*kube.Client), yamlContent[:yamlContentLen])
			// 这样分配的空间还在，不用重复分配，可以利用，将长度设置为0，追加在后面
			yamlContent = yamlContent[:0]
			yamlContentLen = 0
		}
	}
}

func decodeFromYamlContent(client *kube.Client, yamlContent []byte) {
	// parse from yaml
	decode := scheme.Codecs.UniversalDeserializer().Decode
	runtimeObj, gvk, err := decode(yamlContent, nil, nil)
	if err != nil {
		calm_utils.Fatalf("scheme.Codecs.UniversalDeserializer().Decode failed, err:%s", err.Error())
	}

	calm_utils.Debugf("gvk:%s", litter.Sdump(gvk))
	//calm_utils.Debugf("runtimeObj:%#v， type:[%s]", runtimeObj, reflect.TypeOf(runtimeObj).Name())

	cs, _ := client.Factory.KubernetesClientSet()

	// 这里代码可以参考，helm/v3/pkg/kube/wait.go
	// cannot fallthrough in type switch
	// https://stackoverflow.com/questions/11531264/why-isnt-fallthrough-allowed-in-a-type-switch
	switch obj := runtimeObj.(type) {
	case *v1.Service:
		//calm_utils.Debugf("obj type:%s", reflect.TypeOf(obj).Name())
		calm_utils.Debugf("*v1.Service ---> Name:[%s] Namespace[%s]", obj.ObjectMeta.Name, func() string {
			if len(obj.ObjectMeta.Namespace) == 0 {
				return "default"
			}
			return obj.ObjectMeta.Namespace
		}())
	case *appsv1.Deployment:
		// 运行时的状态还是需要去获取
		dp, _ := cs.AppsV1().Deployments(func() string {
			if len(obj.ObjectMeta.Namespace) == 0 {
				return "default"
			}
			return obj.ObjectMeta.Namespace
		}()).Get(obj.ObjectMeta.Name, metav1.GetOptions{})
		calm_utils.Debugf("*appsv1.Deployment ---> Name:[%s] Namespace[%s] status:%#v",
			obj.ObjectMeta.Name, obj.ObjectMeta.Namespace, dp.Status)
	case *appsv1beta1.Deployment:
		// 运行时的状态还是需要去获取
		dp, _ := cs.AppsV1beta1().Deployments(func() string {
			if len(obj.ObjectMeta.Namespace) == 0 {
				return "default"
			}
			return obj.ObjectMeta.Namespace
		}()).Get(obj.ObjectMeta.Name, metav1.GetOptions{})
		calm_utils.Debugf("*appsv1beta1.Deployment ---> Name:[%s] Namespace[%s] status:%#v",
			obj.ObjectMeta.Name, obj.ObjectMeta.Namespace, dp.Status)
	case *appsv1beta2.Deployment:
		// 运行时的状态还是需要去获取
		dp, _ := cs.AppsV1beta2().Deployments(func() string {
			if len(obj.ObjectMeta.Namespace) == 0 {
				return "default"
			}
			return obj.ObjectMeta.Namespace
		}()).Get(obj.ObjectMeta.Name, metav1.GetOptions{})
		calm_utils.Debugf("*appsv1beta2.Deployment ---> Name:[%s] Namespace[%s] status:%#v",
			obj.ObjectMeta.Name, obj.ObjectMeta.Namespace, dp.Status)
	case *extensionsv1beta1.Deployment:
		// 运行时的状态还是需要去获取
		dp, _ := cs.ExtensionsV1beta1().Deployments(func() string {
			if len(obj.ObjectMeta.Namespace) == 0 {
				return "default"
			}
			return obj.ObjectMeta.Namespace
		}()).Get(obj.ObjectMeta.Name, metav1.GetOptions{})
		calm_utils.Debugf("*extensionsv1beta1.Deployment ---> Name:[%s] Namespace[%s] status:%#v",
			obj.ObjectMeta.Name, obj.ObjectMeta.Namespace, dp.Status)
	default:
		objType := reflect.TypeOf(runtimeObj).Name()
		calm_utils.Debugf("obj is not support! objType:[%s] %T", objType, runtimeObj)
	}
}
