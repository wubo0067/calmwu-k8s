/*
 * @Author: calmwu
 * @Date: 2019-09-09 20:15:24
 * @Last Modified by: calmwu
 * @Last Modified time: 2019-09-09 20:17:53
 */

package mysql

import (
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/table"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// job、cronjob的地址存储管理

// SetJobNetInfo 设置job、cronjob的网络信息
func (msm *mysqlStoreMgr) SetJobNetInfo(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	netRegionalID, subNetID, subNetGatewayAddr string) error {
	_, err := msm.dbMgr.Exec(`INSERT INTO tbl_K8SJobNetInfo (k8sresource_id, 
		k8sresource_type, 
		netregional_id, 
		subnet_id,
		subnetgatewayaddr) VALUES (?, ?, ?, ?, ?)`, k8sResourceID, int(k8sResourceType), netRegionalID, subNetID, subNetGatewayAddr)
	if err != nil {
		err = errors.Wrapf(err, "INSERT INTO tbl_K8SJobNetInfo VALUES (%s, %s, %s, %s, %s), Exec failed.", k8sResourceID,
			k8sResourceType.String(), netRegionalID, subNetID, subNetGatewayAddr)
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("SetJobNetInfo k8sResourceID:%s k8sResourceType:%s netRegionalID:%s subNetID:%s subNetGatewayAddr:%s successed.",
		k8sResourceID, k8sResourceType.String(), netRegionalID, subNetID, subNetGatewayAddr)
	return nil
}

// GetJobNetInfo 查询Job、Cronjob的网络信息, 网络域id， 子网id，子网网关地址
func (msm *mysqlStoreMgr) GetJobNetInfo(k8sResourceID string) (string, string, string, error) {
	var k8sJobNetInfo table.TblK8SJobNetInfoS

	err := msm.dbMgr.Get(&k8sJobNetInfo, `SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=? LIMIT 1`, k8sResourceID)
	if err != nil {
		err = errors.Wrapf(err, "SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Get failed", k8sResourceID)
		calm_utils.Error(err.Error())
		return "", "", "", err
	}
	calm_utils.Debugf("SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Get successed.", k8sResourceID)
	return k8sJobNetInfo.NetRegionalID, k8sJobNetInfo.SubNetID, k8sJobNetInfo.SubNetGatewayAddr, nil
}

// DelJobNetInfo 删除Job、Cronjob的网络信息
func (msm *mysqlStoreMgr) DelJobNetInfo(k8sResourceID string) error {
	_, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=? LIMIT 1", k8sResourceID)
	if err != nil {
		err = errors.Wrapf(err, "DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Exec failed.", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Exec successed.", k8sResourceID)
	return nil
}

// BindJobPodWithPortID 绑定job、cronjob的podid和网络地址
func (msm *mysqlStoreMgr) BindJobPodWithPortID(k8sResourceID string, podIP string, portID string, podID string) error {
	return nil
}

// UnbindJobPodWithPortID 解绑job、cronjob的podid和网络地址
func (msm *mysqlStoreMgr) UnbindJobPodWithPortID(k8sResourceID string, podIP string) error {
	return nil
}
