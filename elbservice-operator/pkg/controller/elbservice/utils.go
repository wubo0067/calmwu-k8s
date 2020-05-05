/*
 * @Author: calmwu
 * @Date: 2020-05-04 11:18:54
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-05-04 11:20:04
 */

package elbservice

import (
	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func addFinalizer(meta *metav1.ObjectMeta, finalizer string) {
	if !funk.ContainsString(meta.Finalizers, finalizer) {
		meta.Finalizers = append(meta.Finalizers, finalizer)
	}
}

func removeFinalizer(meta *metav1.ObjectMeta, finalizer string) bool {
	bExists := funk.ContainsString(meta.Finalizers, finalizer)
	if bExists {
		meta.Finalizers = funk.FilterString(meta.Finalizers, func(s string) bool {
			return s != finalizer
		})
	}
	return bExists
}
