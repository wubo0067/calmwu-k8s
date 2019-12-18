/*
 * @Author: calm.wu
 * @Date: 2019-11-29 11:03:09
 * @Last Modified by: calmwu
 * @Last Modified time: 2019-11-30 14:48:51
 */

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/nsp"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func (msm *mysqlStoreMgr) MaintainForceReleaseK8SResourceIPPool(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType) error {

	err := msm.dbSafeExec(context.Background(), func(ctx context.Context) error {

		// 事务
		tx, err := msm.dbMgr.BeginTxx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		if err != nil {
			err := errors.Wrapf(err, "k8sResourceID:%s ResourceType:%s release IPPool begin transaction failed.",
				k8sResourceID, k8sResourceType.String())
			calm_utils.Error(err.Error())
			return err
		}

		var transactionFlag int
		defer func(flag *int) {
			if *flag == 0 {
				calm_utils.Debugf("k8sResourceID:%s ResourceType:%s Commit", k8sResourceID, k8sResourceType.String())
				tx.Commit()
			} else {
				calm_utils.Warnf("k8sResourceID:%s ResourceType:%s Rollback", k8sResourceID, k8sResourceType.String())
				tx.Rollback()
			}
		}(&transactionFlag)

		// 对tbl_K8SResourceIPBind加gap lock，找出要释放的地址资源
		selRows, err := tx.Queryx("SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=?",
			k8sResourceID, k8sResourceType)
		if err != nil {
			err = errors.Wrapf(err, "SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=%s failed",
				k8sResourceID, k8sResourceType.String())
			calm_utils.Error(err.Error())
			transactionFlag = -1
			return err
		}

		// 将地址主动释放给NSP
		var ip, portID string
		for selRows.Next() {
			err = selRows.Scan(&ip, &portID)
			if err == nil {
				calm_utils.Debugf("Return k8s addr{%s---%s} to NSP", ip, portID)
				nsp.NSPMgr.ReleaseAddrResources(portID)
			}
		}
		selRows.Close()

		calm_utils.Debugf("Start --> DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=%s",
			k8sResourceID, k8sResourceType.String())

		// 删除tbl_K8SResourceIPBind表中对应记录
		delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=?",
			k8sResourceID, k8sResourceType)
		if err != nil {
			err = errors.Wrapf(err, "DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=%d failed.",
				k8sResourceID, k8sResourceType)
			calm_utils.Error(err)
		} else {
			delRowCount, _ := delRes.RowsAffected()
			calm_utils.Debugf("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=%s successed. delRowCount:%d",
				k8sResourceID, k8sResourceType.String(), delRowCount)
		}

		calm_utils.Debugf("Start --> DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=%s", k8sResourceID)

		// 删除tbl_K8SResourceIPRecycle表中对应记录
		delRes, err = msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=?",
			k8sResourceID)
		if err != nil {
			err = errors.Wrapf(err, "DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=%s failed.",
				k8sResourceID)
			calm_utils.Error(err)
		} else {
			delRowCount, _ := delRes.RowsAffected()
			calm_utils.Debugf("DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=%s successed. delRowCount:%d",
				k8sResourceID, delRowCount)
		}

		calm_utils.Debugf("Start --> DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s", k8sResourceID)

		// 删除tbl_K8SScaleDownMark表中对应记录
		delRes, err = msm.dbMgr.Exec("DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=?",
			k8sResourceID)
		if err != nil {
			err = errors.Wrapf(err, "DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s failed.",
				k8sResourceID)
			calm_utils.Error(err)
		} else {
			delRowCount, _ := delRes.RowsAffected()
			calm_utils.Debugf("DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s successed. delRowCount:%d",
				k8sResourceID, delRowCount)
		}

		return nil
	})
	return err
}

// MaintainDelCronjobNetInfos cronjob的网络信息释放
func (msm *mysqlStoreMgr) MaintainDelCronjobNetInfos(k8sResourceID string) error {
	cronjobJobLikeName := fmt.Sprintf("%s-%%", k8sResourceID)
	delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=? OR (k8sresource_id LIKE ? AND k8sresource_type=2)",
		k8sResourceID, cronjobJobLikeName)
	if err != nil {
		err = errors.Wrapf(err, "DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s OR (k8sresource_id LIKE %s AND k8sresource_type=2) failed.",
			k8sResourceID, cronjobJobLikeName)
		calm_utils.Error(err)
		return err
	} else {
		delRowCount, _ := delRes.RowsAffected()
		calm_utils.Debugf("DELETE FROM tbl_K8SJobNetInfo WHERE k8sresource_id=%s OR (k8sresource_id LIKE %s AND k8sresource_type=2) successed. delRowCount:%d",
			k8sResourceID, cronjobJobLikeName, delRowCount)
	}
	return nil
}

// MaintainForceReleasePodIP 强制释放pod的ip，归还给nsp，同时删除该记录
func (msm *mysqlStoreMgr) MaintainForceReleasePodIP(k8sResourceID string, bindPodUniqueName string) error {
	var ip, portID string
	row := msm.dbMgr.QueryRow("SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_poduniquename=?",
		k8sResourceID, bindPodUniqueName)
	err := row.Scan(&ip, &portID)
	if err != nil {
		err = errors.Wrapf(err, "SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND bind_poduniquename=%s failed.",
			k8sResourceID, bindPodUniqueName)
		calm_utils.Error(err.Error())
		return err
	}

	calm_utils.Debugf("Return k8s addr{%s---%s} to NSP", ip, portID)
	nsp.NSPMgr.ReleaseAddrResources(portID)

	// 删除该条记录
	delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_poduniquename=? LIMIT 1",
		k8sResourceID, bindPodUniqueName)
	if err != nil {
		err = errors.Wrapf(err, "DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND bind_poduniquename=%s LIMIT 1 failed.",
			k8sResourceID, bindPodUniqueName)
		calm_utils.Error(err)
	} else {
		delRowCount, _ := delRes.RowsAffected()
		calm_utils.Debugf("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND bind_poduniquename=%s LIMIT 1 successed. delRowCount:%d",
			k8sResourceID, bindPodUniqueName, delRowCount)
	}

	return nil
}
