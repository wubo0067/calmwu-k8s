/*
 * @Author: calm.wu
 * @Date: 2019-09-13 11:41:15
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-13 16:17:03
 */

package k8sclient

import (
	"encoding/base64"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"strings"
	"sync"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// K8SClient 集群访问接口
type K8SClient interface {
	//  LoadMultiClusterClient 通过配置数据加载多集群的clientset
	LoadMultiClusterClient(k8sClusterCfgDataLst []config.K8SClusterCfgData) bool

	// GetClientSetByClusterID 获取clientset根据clusterid
	GetClientSetByClusterID(clusterID string) *kubernetes.Clientset

	// 根据cluster-id, namespace，pod-id获取所在节点node的状态
	GetNodeStatus(k8sResourceID string, podID string) error
}

type K8sClientImpl struct {
	multiClusterClient sync.Map
}

var (
	// DefaultK8SClient 默认对象
	DefaultK8SClient K8SClient = &K8sClientImpl{}
)

// LoadMultiClusterClient 通过配置数据加载多集群的clientset
func (kci *K8sClientImpl) LoadMultiClusterClient(k8sClusterCfgDataLst []config.K8SClusterCfgData) bool {
	var loadOk bool = true
	for index := range k8sClusterCfgDataLst {
		k8sClusterCfgData := &k8sClusterCfgDataLst[index]
		// 创建clientset
		kubeCfg, err := base64.StdEncoding.DecodeString(k8sClusterCfgData.KubeCfg)
		if err != nil {
			loadOk = false
			calm_utils.Errorf("cluster:%s base64 DecodeString kubeCfg failed. err:%s", k8sClusterCfgData.K8SClusterID, err.Error())
		} else {
			calm_utils.Debug("cluster:%s \nkubeCfg:%s", k8sClusterCfgData.K8SClusterID, calm_utils.Bytes2String(kubeCfg))
			clientSet, err := NewClientSetByKubeCfgContent(kubeCfg)
			if err != nil {
				calm_utils.Errorf("Load cluster:%s kube config failed. err:%s", k8sClusterCfgData.K8SClusterID, err.Error())
				kci.multiClusterClient.Store(k8sClusterCfgData.K8SClusterID, clientSet)
				loadOk = false
			} else {
				calm_utils.Debugf("Load cluster:%s kube config successed", k8sClusterCfgData.K8SClusterID)
			}
		}
	}
	return loadOk
}

// GetClientSetByClusterID 获取clientset根据clusterid
func (kci *K8sClientImpl) GetClientSetByClusterID(clusterID string) *kubernetes.Clientset {
	value, exist := kci.multiClusterClient.Load(clusterID)
	if exist {
		return value.(*kubernetes.Clientset)
	}
	return nil
}

// GetNodeStatus 根据cluster-id, namespace，pod-id获取所在节点node的状态
func (kci *K8sClientImpl) GetNodeStatus(k8sResourceID string, podID string) error {
	content := k8sResourceID
	pos := strings.IndexByte(k8sResourceID, ':')
	clusterID := k8sResourceID[:pos]
	content = k8sResourceID[pos+1:]
	pos = strings.IndexByte(content, ':')
	namespace := content[:pos]

	calm_utils.Debugf("k8sResourceID:%s clusterID:%s namespace:%s", clusterID, clusterID, namespace)

	clientSet := kci.GetClientSetByClusterID(clusterID)
	if clientSet == nil {
		err := errors.Errorf("clusterID:%s not in config", clusterID)
		calm_utils.Error(err.Error())
		return err
	}

	k8sPod, err := clientSet.CoreV1().Pods(namespace).Get(podID, metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		err := errors.Wrapf(err, "Get Pod:%s Namespace:%s failed.", podID, namespace)
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("clusterID:%s namespace:%s podID:%s Status:%#v", clusterID, namespace, podID, k8sPod.Status)

	nodeName := k8sPod.Spec.NodeName
	k8sNode, err := clientSet.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		err := errors.Wrapf(err, "Get Node:%s failed.", nodeName)
		calm_utils.Error(err.Error())
		return err
	}

	calm_utils.Debugf("node:%s status:%v", nodeName, k8sNode.Status.Conditions)
	return nil
}
