package main

import (
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	//admissionWebHookAnnotationMutateKey = "admission-webhook-calmwu/mutate"
	admissionWebhookAnnotationCreateKey = "admission-webhook-calmwu/create"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

var (
	ignoreSubResouces = map[string]struct{}{
		"status": struct{}{},
	}
)

func serveMutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels, availableAnnotations map[string]string
		objectMeta                            *metav1.ObjectMeta
		resourceNamespace, resourceName       string
	)
	klog.Infof("AdmissionReview for Resource=[%s], SubResource=[%s], Namespace=[%s], Name=[%s], UID=[%s] patchOperation=[%v]"+
		" Kind:%#v, UserInfo=%#v",
		req.Resource.String(), req.SubResource, req.Namespace, req.Name, req.UID, req.Operation, req.Kind, req.UserInfo)

	if req.Name == "webhook-calm-server" ||
		req.Operation == "DELETE" {
		klog.Info("Operation[DELETE] or Name[webhook-calm-server] default allows passed!")
		return nil
	}

	if _, exists := ignoreSubResouces[req.SubResource]; exists {
		klog.Infof("SubResource[%s] default allows passed!", req.SubResource)
		return nil
	}

	switch req.Kind.Kind {
	case "Deployment":
		// 这里要根据不同的req.Kind，Group，Version，Kind来找到具体的对象，Deployment也有多种类型的，要兼容用户的需求。
		var deployment appsv1.Deployment
		// 反序列化
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			klog.Error(err.Error())
			return &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels
		availableAnnotations = objectMeta.GetAnnotations()

		// 判断etcd中是否有该deployment的地址集合，没有就创建，同时比较Replicas的数量。
		// 通过这个去nsp获取pod地址集合
		klog.Infof("deployment.spec.replicas:%d", *deployment.Spec.Replicas)
		klog.Infof("deployment:%#v", deployment)
	case "StatefulSet":
		var statefulset appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &statefulset); err != nil {
			klog.Error(err.Error())
			return &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = statefulset.Name, statefulset.Namespace, &statefulset.ObjectMeta
		availableLabels = statefulset.Labels
		availableAnnotations = objectMeta.GetAnnotations()
		klog.Infof("statefulset.spec.replicas:%d", *statefulset.Spec.Replicas)
		klog.Infof("statefulset:%#v", statefulset)
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(req.Object.Raw, &service); err != nil {
			klog.Error(err.Error())
			return &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels
		availableAnnotations = objectMeta.GetAnnotations()
	case "Scale":
		var scale extensionv1beta1.Scale
		if err := json.Unmarshal(req.Object.Raw, &scale); err != nil {
			klog.Error(err.Error())
			return &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		klog.Infof("scale:%#v", scale)
		resourceName, resourceNamespace, objectMeta = scale.Name, scale.Namespace, &scale.ObjectMeta
		availableLabels = scale.Labels
		availableAnnotations = objectMeta.GetAnnotations()
	default:
		klog.Infof("this Kind:%s resource using default process", req.Kind.Kind)
		return nil
	}

	klog.Infof("resourceName:%s", resourceName)
	klog.Infof("resourceNamespace:%s", resourceNamespace)
	klog.Infof("available Labels:%#v", availableLabels)
	klog.Infof("available Annotations:%#v", availableAnnotations)

	// add patch for /metadata/annotations
	value := fmt.Sprintf("calmwu-%s", time.Now().Format("2006-01-02 15:04:05 -0700"))
	var patches []patchOperation

	// 判断这个key是否已经存在
	if _, exist := availableAnnotations[admissionWebhookAnnotationCreateKey]; exist {
		patches = append(patches, patchOperation{
			Op:    "replace",
			Path:  "/metadata/annotations/" + admissionWebhookAnnotationCreateKey, // 如果修改，直接将key拼接到path中。
			Value: value,
		})
	} else {
		patches = append(patches, patchOperation{
			Op:   "add",
			Path: "/metadata/annotations",
			Value: map[string]string{
				admissionWebhookAnnotationCreateKey: value, // 第一次要将key作为map的key
			},
		})
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		errStr := fmt.Sprintf("json marsh patchOperation failed, reason:%s", err.Error())
		klog.Error(errStr)
		return &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: errStr,
			},
		}
	}

	klog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
