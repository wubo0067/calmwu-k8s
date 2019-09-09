/*
 * @Author: calmwu
 * @Date: 2019-09-09 20:15:24
 * @Last Modified by: calmwu
 * @Last Modified time: 2019-09-09 20:17:53
 */

package mysql

import proto "pci-ipresmgr/api/proto_json"

// job、cronjob的地址存储管理

// SetJobNetInfo 设置job、cronjob的网络信息
func (msm *mysqlStoreMgr) SetJobNetInfo(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	netRegionalID, subNetID, subNetGatewayAddr string) error {
	return nil
}

// GetJobNetInfo 查询Job、Cronjob的网络信息
func (msm *mysqlStoreMgr) GetJobNetInfo(k8sResourceID string) (string, string, string, error) {

}

// BindJobPodWithPortID 绑定job、cronjob的podid和网络地址
func (msm *mysqlStoreMgr) BindJobPodWithPortID(k8sResourceID string, podIP string, portID string, podID string) error {
	return nil
}

// UnbindJobPodWithPortID 解绑job、cronjob的podid和网络地址
func (msm *mysqlStoreMgr) UnbindJobPodWithPortID(k8sResourceID string, podIP string) error {
	return nil
}
