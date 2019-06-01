/*
 * @Author: calm.wu
 * @Date: 2019-05-24 14:55:51
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-24 14:58:48
 */

package main

import (
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func serveMutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels, availableAnnotations map[string]string
		objectMeta                            *metav1.ObjectMeta
		resourceNamespace, resourceName       string
	)

	klog.Infof("AdmissionReview for Resource=%s, Namespace=%s, Name=%s, UID=%s patchOperation=%v, UserInfo=%v",
		req.Resource.String(), req.Namespace, req.Name, req.UID, req.Operation, req.UserInfo)

	switch req.Kind.Kind {

	}

	return nil
}
