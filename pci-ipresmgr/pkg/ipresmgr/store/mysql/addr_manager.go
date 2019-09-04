/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:52:15
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 17:28:35
 */

package mysql

import (
	"context"
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/store"
	"strings"
	"time"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// GetAddrCountByK8SResourceID 根据资源id名，获取k8s资源对应的地址数量
func (msm *mysqlStoreMgr) GetAddrCountByK8SResourceID(k8sReousrceID string) (int, error) {
	return 0, nil
}

// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
func (msm *mysqlStoreMgr) SetAddrInfosToK8SResourceID(k8sReousrceID string, k8sResourceType proto.K8SApiResourceKindType,
	k8sAddrInfos []*store.K8SAddrInfo) error {
	// 插入tbl_K8SResourceIPBind表
	return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		tx := msm.dbMgr.MustBegin()

		for index, k8sAddrInfo := range k8sAddrInfos {
			calm_utils.Debugf("%d k8sReousrceID[%s] k8sResourceType[%s] k8sAddrInfo:%+v", index, k8sReousrceID,
				k8sResourceType.String(), k8sAddrInfo)

			insRes := tx.MustExec(`INSERT INTO tbl_K8SResourceIPBind 
			k8sresource_id, 
			k8sresource_type, 
			ip, 
			mac,
			netregional_id,
			subnet_id,
			port_id,
			subnetgatewayaddr,
			alloc_time,
			is_bind,
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				k8sReousrceID,
				int(k8sResourceType),
				k8sAddrInfo.IP,
				k8sAddrInfo.MacAddr,
				k8sAddrInfo.NetRegionalID,
				k8sAddrInfo.SubNetID,
				k8sAddrInfo.PortID,
				k8sAddrInfo.SubNetGatewayAddr,
				time.Now(),
				0,
			)

			_, err := insRes.RowsAffected()
			if err != nil {
				err = errors.Wrapf(err, "insert k8sReousrceID[%s] k8sResourceType[%s] addrinfo failed", k8sReousrceID, k8sResourceType.String())
				calm_utils.Error(err.Error())
				tx.Rollback()
				return err
			}
		}

		tx.Commit()
		return nil
	})

	return nil
}

// GetAddrInfoByK8SResourceID 获取一个地址信息
func (msm *mysqlStoreMgr) GetAddrInfoByK8SResourceID(k8sReousrceID string) *store.K8SAddrInfo {
	k8sAddrInfo := new(store.K8SAddrInfo)
	return k8sAddrInfo
}

// CheckRecycledResources 检查对应资源是否存在，bool = true存在，int=副本数量
func (msm *mysqlStoreMgr) CheckRecycledResources(k8sReousrceID string) (bool, int, error) {
	// 在tbl_K8SResourceIPRecycle查询
	var replicas int
	err := msm.dbMgr.Get(&replicas, "SELECT replicas, unbind_count FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=?", k8sReousrceID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// 地址资源不在待回收表中，
			calm_utils.Infof("%s not in tbl_K8SResourceIPRecycle, must be call NSP get addr resource.", k8sReousrceID)
			return false, 0, nil
		}
		return false, 0, errors.Wrapf(err, "CheckRecycledResources %s failed.", k8sReousrceID)
	}

	// 说明存在待回收表，删除该记录，恢复
	delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=?", k8sReousrceID)
	if err != nil {
		err = errors.Wrapf(err, "DELETE tbl_K8SResourceIPRecycle by %s failed", k8sReousrceID)
		calm_utils.Fatalf(err.Error())
		return false, 0, err
	}

	delRows, _ := delRes.RowsAffected()
	if delRows != 1 {
		err = errors.Wrapf(err, "DELETE tbl_K8SResourceIPRecycle by %s delRows:%d is incorrect", k8sReousrceID, delRows)
		calm_utils.Error(err)
		return false, 0, err
	} else {
		calm_utils.Infof("Restored resources:%s to be recovered", k8sReousrceID)
	}

	return true, replicas, nil
}
