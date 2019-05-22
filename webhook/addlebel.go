/*
 * @Author: calm.wu
 * @Date: 2019-05-22 11:04:46
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-22 15:35:45
 */

package main

import (
	"encoding/json"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	addFirstLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels", "value": {"added-label": "yes", "add-by": "calmwu"}}
     ]`
	addAdditionalLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels/added-label", "value": "yes" }
     ]`
)

func addLabel(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(0).Info("calling add-label")

	obj := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	klog.Infof("addLabel obj:%#v", obj)

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if len(obj.ObjectMeta.Labels) == 0 {
		reviewResponse.Patch = []byte(addFirstLabelPatch)
	} else {
		reviewResponse.Patch = []byte(addAdditionalLabelPatch)
	}
	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse	

	return nil
}