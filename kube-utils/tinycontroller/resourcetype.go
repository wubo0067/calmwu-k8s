/*
 * @Author: calm.wu
 * @Date: 2020-09-02 15:24:29
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-02 15:28:56
 */

package tinycontroller

type ResourceType string

const (
	//
	Deployment            ResourceType = "deployment"
	ReplicationController ResourceType = "rc"
	ReplicaSet            ResourceType = "rs"
	DaemonSet             ResourceType = "ds"
	Service               ResourceType = "svc"
	Pod                   ResourceType = "po"
	Job                   ResourceType = "job"
	Node                  ResourceType = "node"
	ClusterRole           ResourceType = "clusterrole"
	ServiceAccount        ResourceType = "sa"
	PersistentVolume      ResourceType = "pv"
	Namespace             ResourceType = "ns"
	Secret                ResourceType = "secret"
	ConfigMap             ResourceType = "configmap"
	Ingress               ResourceType = "ing"
	EndPoints             ResourceType = "ep"
)
