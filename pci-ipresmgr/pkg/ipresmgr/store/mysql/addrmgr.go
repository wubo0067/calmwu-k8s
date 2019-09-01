/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:52:15
 * @Last Modified by:   calm.wu
 * @Last Modified time: 2019-09-01 09:52:15
 */

package mysql

import "pci-ipresmgr/pkg/ipresmgr/store"

// GetAddrCountByK8SResourceID 根据资源id名，获取k8s资源对应的地址数量
func (msm *mysqlStoreMgr) GetAddrCountByK8SResourceID(K8SReousrceID string) (int, error) {
	return 0, nil
}

// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
func (msm *mysqlStoreMgr) SetAddrInfosToK8SResourceID(K8SResourceID string, k8sAddrInfos []*store.K8SAddrInfo) error {
	return nil
}

// GetAddrInfoByK8SResourceID 获取一个地址信息
func (msm *mysqlStoreMgr) GetAddrInfoByK8SResourceID(K8SReousrceID string) *store.K8SAddrInfo {
	k8sAddrInfo := new(store.K8SAddrInfo)
	return k8sAddrInfo
}
