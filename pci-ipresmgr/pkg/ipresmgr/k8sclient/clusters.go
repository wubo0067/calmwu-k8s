/*
 * @Author: calm.wu
 * @Date: 2019-09-13 11:41:15
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-13 15:13:11
 */

package k8sclient

import (
	"pci-ipresmgr/pkg/ipresmgr/config"
	"sync"
)

var (
	multiClusterClient sync.Map
)

// LoadMultiClusterClient 通过配置数据加载多集群的clientset
func LoadMultiClusterClient(k8sClusterCfgDataLst []config.K8SClusterCfgData) error {
	return nil
}
