/*
 * @Author: calm.wu
 * @Date: 2020-02-07 17:11:54
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-02-07 17:18:03
 */

package resources

import (
	k8sv1alpha1 "calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPodForCR创建一个pod资源
func NewPodForCR(cr *k8sv1alpha1.ELBService) *corev1.Pod {
	labels := map[string]string{
		"app":     cr.Name,
		"version": "v0.1",
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-pod",
			Namespace:    cr.Namespace,
			Labels:       labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					Command:         []string{"sleep", "3600"},
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
		},
	}
}
