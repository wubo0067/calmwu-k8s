/*
 * @Author: calm.wu
 * @Date: 2020-09-07 11:28:40
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-07 11:31:14
 */

package tinycontroller

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetObjectMetaData returns metadata of a given k8s object
func GetObjectMetaData(obj interface{}) (objectMeta metav1.ObjectMeta) {

	switch object := obj.(type) {
	case *appsv1.Deployment:
		objectMeta = object.ObjectMeta
	case *apiv1.ReplicationController:
		objectMeta = object.ObjectMeta
	case *appsv1.ReplicaSet:
		objectMeta = object.ObjectMeta
	case *appsv1.DaemonSet:
		objectMeta = object.ObjectMeta
	case *apiv1.Service:
		objectMeta = object.ObjectMeta
	case *apiv1.Pod:
		objectMeta = object.ObjectMeta
	case *batchv1.Job:
		objectMeta = object.ObjectMeta
	case *apiv1.PersistentVolume:
		objectMeta = object.ObjectMeta
	case *apiv1.Namespace:
		objectMeta = object.ObjectMeta
	case *apiv1.Secret:
		objectMeta = object.ObjectMeta
	case *extv1beta1.Ingress:
		objectMeta = object.ObjectMeta
	case *apiv1.Node:
		objectMeta = object.ObjectMeta
	case *rbacv1beta1.ClusterRole:
		objectMeta = object.ObjectMeta
	case *apiv1.ServiceAccount:
		objectMeta = object.ObjectMeta
	case *apiv1.Event:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}
