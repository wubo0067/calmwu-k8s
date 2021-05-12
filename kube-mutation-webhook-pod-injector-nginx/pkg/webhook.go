/*
 * @Author: CALM.WU
 * @Date: 2021-04-29 15:28:02
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-29 17:51:49
 */

package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"istio.io/istio/pkg/kube"
	admissionregistrationv1 "k8s.io/api/admission/v1"
	admissionregistrationv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// 空schema
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

const (
	_admissionWebhookAnnotationInjectKey = "nginx-injector-pod-webhook/inject"
	_admissionWebhookAnnotationStatusKey = "nginx-injector-pod-webhook/status"
)

func init() {
	// 注册类型的schema
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1.AddToScheme(runtimeScheme)
}

func toAdmissionResponse(err error) *kube.AdmissionResponse {
	return &kube.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// HandleInject call by apiserver
func HandleInject(w http.ResponseWriter, r *http.Request) {
	var body []byte

	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			glog.Infof("read request body successed! body size: %d", len(data))
			body = data
		}
	}

	if len(body) == 0 {
		glog.Error("request body is empty")
		http.Error(w, "body empty", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	path := ""
	if r.URL != nil {
		path = r.URL.Path
	}

	// 响应对象
	var admissionResponse *kube.AdmissionResponse
	// 请求结构为AdmissionReview
	var reqAR *kube.AdmissionReview
	var obj runtime.Object

	if out, _, err := deserializer.Decode(body, nil, obj); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = toAdmissionResponse(err)
	} else {
		// 默认通过
		// runtime.Object ===> v1beta1.AdmissionResponse / v1.AdmissionReview 支持两种类型，统一转换为适配类型
		reqAR, err = kube.AdmissionReviewKubeToAdapter(out)
		if err != nil {
			// 类型转换失败
			admissionResponse = toAdmissionResponse(err)
		} else {
			glog.Infof("Inject AdmissionReview kind:%s apiVersion:%s uuid:%s for path:%s",
				reqAR.Kind, reqAR.APIVersion, reqAR.Request.UID, path)
			// 这里执行inject
			admissionResponse = inject(reqAR, path)
		}
	}

	// 返回对象
	resAR := kube.AdmissionReview{}
	resAR.Response = admissionResponse

	var apiVersion string

	if reqAR != nil {
		apiVersion = reqAR.APIVersion
		resAR.TypeMeta = reqAR.TypeMeta
		if resAR.Response != nil {
			if reqAR.Request != nil {
				resAR.Response.UID = reqAR.Request.UID
			}
		}
	}

	resRuntimeObj := kube.AdmissionReviewAdapterToKube(&resAR, apiVersion)

	resp, err := json.Marshal(resRuntimeObj)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// inject start inject nginx container into pod
func inject(ar *kube.AdmissionReview, path string) *kube.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod

	// 从raw数据中解析出pod对象
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		err := errors.Wrap(err, "Could not unmarshal raw object!")
		glog.Error(err.Error())
		return toAdmissionResponse(err)
	}

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, pod.Name, req.UID, req.Operation, req.UserInfo)

	// 根据namespace和pod的annotation判断是否该注入sidecar
	if !injectRequired(&pod.Spec, &pod.ObjectMeta) {
		glog.Infof("Skipping inject for %s/%s due to policy check", pod.Namespace, pod.Name)
		return &kube.AdmissionResponse{
			Allowed: true,
		}
	}

	annotations := map[string]string{_admissionWebhookAnnotationStatusKey: "injected"}
	patchBytes, err := createPatch(&pod, annotations)
	if err != nil {
		err = errors.Wrap(err, "Create pod patch failed.")
		glog.Error(err.Error())
		return toAdmissionResponse(err)
	}

	glog.Infof("AdmissionResponse: patch=%s\n", patchBytes)
	return &kube.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *string {
			pt := "JSONPatch"
			return &pt
		}(),
	}
}

func injectRequired(podSpec *corev1.PodSpec, metadata *metav1.ObjectMeta) bool {
	if podSpec.HostNetwork {
		return false
	}

	if isIgnoreNamespace(metadata.Namespace) {
		glog.Infof("Pod %s namespace %s don't inject", metadata.Name, metadata.Namespace)
		return false
	}

	var inject bool
	// 检查是否有注入的annotation
	annotations := metadata.GetAnnotations()
	switch strings.ToLower(annotations[_admissionWebhookAnnotationInjectKey]) {
	case "y", "yes", "true", "on":
		inject = true
	default:
		inject = false
	}

	glog.Infof("Pod %s/%s: inject required:%v", metadata.Namespace, metadata.Name, inject)
	return inject
}

func createPatch(pod *corev1.Pod, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation

	patch = append(patch, addContainer(pod.Spec.Containers, getSidecarConfig().Containers, "/spec/containers")...)
	patch = append(patch, addVolume(pod.Spec.Volumes, getSidecarConfig().Volumes, "/spec/volumes")...)
	patch = append(patch, updateAnnotaion(pod.Annotations, annotations)...)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(patch)
}

func addContainer(baseContainers, addContainers []corev1.Container, basePath string) (patch []patchOperation) {
	// 每个容器对象加入
	firstContainer := len(baseContainers) == 0
	var val interface{}

	glog.Infof("firstContainer is %v", firstContainer)

	for _, addContainer := range addContainers {
		val = addContainer
		path := basePath

		if firstContainer {
			// 如果没有容器对象，val必须是个sliace对象
			firstContainer = false
			val = []corev1.Container{addContainer}
		} else {
			path = path + "/-"
		}

		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: val,
		})
	}
	return
}

func addVolume(baseVolumes, addVolumes []corev1.Volume, basePath string) (patch []patchOperation) {
	// 加入卷对象
	// firstVolume := len(baseVolumes) == 0
	// var val interface{}

	volumes := append(baseVolumes, addVolumes...)

	glog.Infof("volumes is %#v", volumes)

	patch = append(patch, patchOperation{
		Op:    "replace",
		Path:  basePath,
		Value: volumes,
	})

	// for _, addVolume := range addVolumes {
	// 	val = addVolume
	// 	path := basePath

	// 	if firstVolume {
	// 		firstVolume = false
	// 		val = []corev1.Volume{addVolume}
	// 	} else {
	// 		path = path + "/-"
	// 	}

	// 	patch = append(patch, patchOperation{
	// 		Op:    "add",
	// 		Path:  path,
	// 		Value: val,
	// 	})
	// }
	return
}

func updateAnnotaion(baseAnnotations, addAnnotations map[string]string) (patch []patchOperation) {
	if baseAnnotations == nil {
		baseAnnotations = map[string]string{}
	}

	for key, value := range addAnnotations {
		if baseAnnotations[key] == "" {
			// 之前不存在annotation，或annotation中该key不准在
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return
}
