/*
 * @Author: calm.wu
 * @Date: 2020-02-07 11:42:48
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-02-07 17:18:12
 */

// Package resources for Implement the resources included in elbservice
package resources

import (
	"fmt"

	k8sv1alpha1 "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetVIPServiceName(cr *k8sv1alpha1.ELBService) string {
	return fmt.Sprintf("%s-vipsvc", cr.GetName())
}

func GetEPServiceName(cr *k8sv1alpha1.ELBService) string {
	return fmt.Sprintf("%s-endsvc", cr.GetName())
}

// NewVIPServiceForCR 创建一个无头service，没有selector
func NewVIPServiceForCR(cr *k8sv1alpha1.ELBService) *corev1.Service {
	vipService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetVIPServiceName(cr),
			Namespace: cr.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       cr.Spec.Listener.Name,
					Protocol:   corev1.Protocol(cr.Spec.Listener.Protocol),
					Port:       cr.Spec.Listener.Port,
					TargetPort: intstr.FromInt(int(cr.Spec.Listener.Port)),
				},
			},
		},
	}
	return vipService
}

func NewEPServiceForCR(cr *k8sv1alpha1.ELBService) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetEPServiceName(cr),
			Namespace: cr.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None", // 这个是headless service
			Selector:  cr.Spec.Selector,
		},
	}
}
