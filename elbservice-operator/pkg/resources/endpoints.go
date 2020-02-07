/*
 * @Author: calm.wu
 * @Date: 2020-02-07 11:43:58
 * @Last Modified by:   calm.wu
 * @Last Modified time: 2020-02-07 11:43:58
 */

package resources

import (
	k8sv1alpha1 "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewVIPEndpointForCR创建一个endpoints资源
func NewVIPEndpointForCR(cr *k8sv1alpha1.ELBService) *corev1.Endpoints {
	vipEndpoints := &corev1.Endpoints{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Endpoints",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetVIPServiceName(cr),
			Namespace: cr.GetNamespace(),
		},
		Subsets: []corev1.EndpointSubset{
			corev1.EndpointSubset{
				Addresses: []corev1.EndpointAddress{
					corev1.EndpointAddress{
						IP: cr.Spec.Listener.VIP, // 设置elb的ip地址
					},
				},
				Ports: []corev1.EndpointPort{
					corev1.EndpointPort{
						Port:     cr.Spec.Listener.Port,
						Protocol: corev1.Protocol(cr.Spec.Listener.Protocol),
					},
				},
			},
		},
	}

	return vipEndpoints
}
