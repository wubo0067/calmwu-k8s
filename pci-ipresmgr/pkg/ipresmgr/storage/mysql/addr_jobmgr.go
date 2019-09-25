/*
 * @Author: calmwu
 * @Date: 2019-09-09 20:15:24
 * @Last Modified by: calmwu
 * @Last Modified time: 2019-09-09 20:17:53
 */

package mysql

import (
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/nsp"
	"pci-ipresmgr/table"
	"time"

	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// job、cronjob的地址存储管理

// SetJobNetInfo 设置job、cronjob的网络信息
func (msm *mysqlStoreMgr) SetJobNetInfo(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	netRegionalID, subNetID, subNetGatewayAddr string, subNetCIDR string) error {
	now := time.Now()
	_, err := msm.dbMgr.Exec(`INSERT INTO tbl_K8SJobNetInfo (k8sresource_id, 
		k8sresource_type, 
		netregional_id, 
		subnet_id,
		subnetgatewayaddr,
		subnetcidr,
		create_time) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		k8sResourceID, int(k8sResourceType), netRegionalID, subNetID, subNetGatewayAddr, subNetCIDR, now)
	if err != nil {
		err = errors.Wrapf(err, "INSERT INTO tbl_K8SJobNetInfo VALUES (%s, %s, %s, %s, %s, %s, %s), Exec failed.", k8sResourceID,
			k8sResourceType.String(), netRegionalID, subNetID, subNetGatewayAddr, subNetCIDR, now.String())
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("SetJobNetInfo k8sResourceID:%s k8sResourceType:%s netRegionalID:%s subNetID:%s subNetGatewayAddr:%s subNetCIDR:%s createTime:%s successed.",
		k8sResourceID, k8sResourceType.String(), netRegionalID, subNetID, subNetGatewayAddr, subNetCIDR, now.String())
	return nil
}

// GetJobNetInfo 查询Job、Cronjob的网络信息, 网络域id， 子网id，子网网关地址
func (msm *mysqlStoreMgr) GetJobNetInfo(k8sResourceID string) (string, string, string, string, error) {
	var k8sJobNetInfo table.TblK8SJobNetInfoS

	err := msm.dbMgr.Get(&k8sJobNetInfo, `SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=? LIMIT 1`, k8sResourceID)
	if err != nil {
		err = errors.Wrapf(err, "SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Get failed", k8sResourceID)
		calm_utils.Error(err.Error())
		return "", "", "", "", err
	}
	calm_utils.Debugf("SELECT * FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s LIMIT 1, Get successed.", k8sResourceID)
	return k8sJobNetInfo.NetRegionalID, k8sJobNetInfo.SubNetID, k8sJobNetInfo.SubNetGatewayAddr, k8sJobNetInfo.SubNetCIDR, nil
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
	now := time.Now()
	_, err := msm.dbMgr.Exec(`INSERT INTO tbl_K8SJobIPBind (k8sresource_id, 
		ip, 
		bind_podid, 
		port_id,
		bind_time ) VALUES (?, ?, ?, ?, ?)`, k8sResourceID, podIP, podID, portID, now)
	if err != nil {
		err = errors.Wrapf(err, "INSERT INTO tbl_K8SJobIPBind VALUES (%s, %s, %s, %s, %s), Exec failed.", k8sResourceID,
			podIP, podID, portID, now.String())
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("INSERT INTO tbl_K8SJobIPBind VALUES (%s, %s, %s, %s, %s), Exec successed.", k8sResourceID,
		podIP, podID, portID, now.String())
	return nil
}

// UnbindJobPodWithPortID 解绑job、cronjob的podid和网络地址
func (msm *mysqlStoreMgr) UnbindJobPodWithPortID(k8sResourceID string, podID string) error {
	var k8sJobIPBind table.TblK8SJobIPBindS

	err := msm.dbMgr.Get(&k8sJobIPBind, `SELECT * FROM tbl_K8SJobIPBind WHERE k8sresource_id=? AND bind_podid=? LIMIT 1`,
		k8sResourceID, podID)
	if err != nil {
		err = errors.Wrapf(err, "SELECT * FROM tbl_K8SJobIPBind WHERE k8sresource_id=%s AND bind_podid=%s LIMIT 1, Get failed",
			k8sResourceID, podID)
		calm_utils.Error(err.Error())
		return err
	}
	calm_utils.Debugf("SELECT * FROM tbl_K8SJobIPBind WHERE k8sresource_id=%s AND bind_podid=%s LIMIT 1, Get successed. k8sJobIPBind:%s",
		k8sResourceID, podID, litter.Sdump(&k8sJobIPBind))

	// 删除
	_, err = msm.dbMgr.Exec("DELETE FROM tbl_K8SJobIPBind WHERE k8sresource_id=? AND bind_podid=? LIMIT 1",
		k8sResourceID, podID)
	if err != nil {
		err = errors.Wrapf(err, "DELETE FROM tbl_K8SJobIPBind WHERE k8sresource_id=%s AND bind_podid=%s LIMIT 1, Exec failed.",
			k8sResourceID, podID)
	}
	calm_utils.Debugf("DELETE FROM tbl_K8SJobIPBind WHERE k8sresource_id=%s AND bind_podid=%s LIMIT 1, Exec successed.",
		k8sResourceID, podID)
	calm_utils.Debugf("Return k8s addr{%s---%s} to NSP", k8sJobIPBind.IP, k8sJobIPBind.PortID)
	nsp.NSPMgr.ReleaseAddrResources(k8sJobIPBind.PortID)
	return nil
}
