/*
 * @Author: calm.wu
 * @Date: 2019-05-24 14:55:51
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-24 14:58:48
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	admissionWebHookAnnotationMutateKey = "admission-webhook-calmwu/mutate"
	admissionWebhookAnnotationCreateKey = "admission-webhook-calmwu/create"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func serveMutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels, availableAnnotations map[string]string
		objectMeta                            *metav1.ObjectMeta
		resourceNamespace, resourceName       string
	)

	klog.Infof("AdmissionReview for Resource=%s, Namespace=%s, Name=%s, UID=%s patchOperation=%v, UserInfo=%v",
		req.Resource.String(), req.Namespace, req.Name, req.UID, req.Operation, req.UserInfo)

	if req.Operation == "DELETE" {
		klog.Info("Operation[DELETE] default passed!")
		return &v1beta1.AdmissionResponse{
			UID:     req.UID,
			Allowed: true,
		}
	}

	switch req.Kind.Kind {
	case "Deployment":
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
		klog.Infof("deployment:%#v", deployment)
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
			Op:   "replace",
			Path: "/metadata/annotations/" + admissionWebhookAnnotationCreateKey,
			Value: map[string]string{
				"CreateBy": value,
			},
		})
	} else {
		patches = append(patches, patchOperation{
			Op:   "add",
			Path: "/metadata/annotations/" + admissionWebhookAnnotationCreateKey,
			Value: map[string]string{
				"CreateBy": value,
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
