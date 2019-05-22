/*
 * @Author: calm.wu 
 * @Date: 2019-05-22 11:06:49 
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-22 14:40:46
 */

package main

import (
	"k8s.io/api/admission/v1beta1"
)

func admitPods(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	return nil
}