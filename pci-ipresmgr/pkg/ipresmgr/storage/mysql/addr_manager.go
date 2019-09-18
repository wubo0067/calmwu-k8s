/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:52:15
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:53:59
 */

package mysql

import (
	"context"
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/nsp"
	"pci-ipresmgr/table"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// GetAddrCountByK8SResourceID 根据资源id名，获取k8s资源对应的地址数量
func (msm *mysqlStoreMgr) GetAddrCountByK8SResourceID(k8sResourceID string) (int, error) {
	return 0, nil
}

// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
func (msm *mysqlStoreMgr) SetAddrInfosToK8SResourceID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	k8sAddrInfos []*proto.K8SAddrInfo) error {
	// 插入tbl_K8SResourceIPBind表
	return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		tx := msm.dbMgr.MustBegin()

		for index, k8sAddrInfo := range k8sAddrInfos {
			calm_utils.Debugf("%d k8sResourceID[%s] k8sResourceType[%s] k8sAddrInfo:%+v", index, k8sResourceID,
				k8sResourceType.String(), k8sAddrInfo)

			insRes := tx.MustExec(`INSERT INTO tbl_K8SResourceIPBind 
			(k8sresource_id, k8sresource_type, ip, mac, netregional_id, subnet_id, port_id, subnetgatewayaddr, alloc_time, is_bind) VALUES 
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				k8sResourceID,
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
				err = errors.Wrapf(err, "insert k8sResourceID[%s] k8sResourceType[%s] addrinfo failed, Rollback.", k8sResourceID, k8sResourceType.String())
				calm_utils.Error(err.Error())
				tx.Rollback()
				return err
			}
		}

		tx.Commit()
		return nil
	})
}

// BindAddrInfoWithK8SPodID 获取一个地址信息，和k8s资源绑定
func (msm *mysqlStoreMgr) BindAddrInfoWithK8SPodID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	bindPodID string) (*proto.K8SAddrInfo, error) {

	var k8sAddrInfo *proto.K8SAddrInfo

	err := msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		if k8sResourceType == proto.K8SApiResourceKindDeployment {

			tx, err := msm.dbMgr.Begin()
			if err != nil {
				err := errors.Wrapf(err, "k8sResourceID:%s bindPod:%s begin transaction failed.", k8sResourceID, bindPodID)
				calm_utils.Error(err.Error())
				return err
			}

			var transactionFlag int
			defer func(flag *int) {
				if *flag == 0 {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Commit", k8sResourceID, bindPodID)
					tx.Commit()
				} else {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Rollback", k8sResourceID, bindPodID)
					tx.Rollback()
				}
			}(&transactionFlag)

			selRow := tx.QueryRow("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 LIMIT 1 FOR UPDATE",
				k8sResourceID, k8sResourceType)
			if selRow == nil {
				err = errors.Errorf("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=%s AND is_bind=0 LIMIT 1, QueryRow return Nil",
					k8sResourceID, k8sResourceType)
				transactionFlag = -1
				return err
			}

			var k8sAddrBindInfo table.TblK8SResourceIPBindS
			err = selRow.Scan(&k8sAddrBindInfo.K8SResourceID, &k8sAddrBindInfo.K8SResourceType, &k8sAddrBindInfo.IP,
				&k8sAddrBindInfo.MacAddr, &k8sAddrBindInfo.NetRegionalID, &k8sAddrBindInfo.SubNetID, &k8sAddrBindInfo.PortID,
				&k8sAddrBindInfo.SubNetGatewayAddr, &k8sAddrBindInfo.AllocTime, &k8sAddrBindInfo.IsBind, &k8sAddrBindInfo.BindPodID,
				&k8sAddrBindInfo.BindTime)

			if err != nil {
				err = errors.Wrapf(err, "k8sResourceID:%s bindPod:%s Scan failed.", k8sResourceID, bindPodID)
				calm_utils.Error(err.Error())
				transactionFlag = -1
				return err
			}

			calm_utils.Debugf("k8sResourceID:%s bindPod:%s k8sAddrBindInfo:%s", k8sResourceID, bindPodID, litter.Sdump(&k8sAddrBindInfo))

			currTime := time.Now()
			updateRes, err := tx.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=?, bind_time=? WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 AND port_id=?",
				bindPodID, currTime, k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.PortID)
			if err != nil {
				err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=%s, bind_time=%s WHERE k8sresource_id=%s AND k8sresource_type=%d AND is_bind=0 AND port_id=%s. tx Exec UPDATE failed.",
					bindPodID, currTime.String(), k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.PortID)
				transactionFlag = -1
				return err
			}

			updateRowCount, _ := updateRes.RowsAffected()
			calm_utils.Debugf("k8sResourceID:%s bindPod:%s updateRowCount:%d\n", k8sResourceID, bindPodID, updateRowCount)

			k8sAddrInfo = new(proto.K8SAddrInfo)
			k8sAddrInfo.IP = k8sAddrBindInfo.IP
			k8sAddrInfo.MacAddr = k8sAddrBindInfo.MacAddr
			k8sAddrInfo.NetRegionalID = k8sAddrBindInfo.NetRegionalID
			k8sAddrInfo.SubNetID = k8sAddrBindInfo.SubNetID
			k8sAddrInfo.SubNetGatewayAddr = k8sAddrBindInfo.SubNetGatewayAddr
			k8sAddrInfo.PortID = k8sAddrBindInfo.PortID

			// 放弃使用悲观锁，使用乐观锁CAS方法。获取，设置，失败要重试
			// step 1，获取重试次数
			// var replicas int
			// err := msm.dbMgr.Get(&replicas, "SELECT count(*) from tbl_K8SResourceIPBind WHERE k8sresource_id=?", k8sResourceID)
			// if err != nil {
			// 	calm_utils.Errorf("SELECT count(*) from tbl_K8SResourceIPBind WHERE k8sresource_id=%s failed. err:%s",
			// 		k8sResourceID, err.Error())
			// 	return nil
			// }

			// if replicas == 0 {
			// 	calm_utils.Errorf("SELECT count(*) from tbl_K8SResourceIPBind WHERE k8sresource_id=%s failed. replicas is Zero",
			// 		k8sResourceID)
			// 	return nil
			// }

			// for replicas > 0 {
			// 	// 查表，获取地址
			// 	var k8sAddrBindInfo table.TblK8SResourceIPBindS
			// 	err := msm.dbMgr.Get(&k8sAddrBindInfo, "SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 LIMIT 1",
			// 		k8sResourceID, int(k8sResourceType))
			// 	if err != nil {
			// 		calm_utils.Errorf("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 failed. err:%s, do try:%d",
			// 			k8sResourceID, err.Error(), replicas)
			// 	} else {
			// 		// 修改表
			// 		currTime := time.Now()
			// 		updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=?, bind_time=? WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 AND ip=?",
			// 			bindPodID, currTime, k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.IP)
			// 		if err != nil {
			// 			calm_utils.Errorf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s failed. err:%s. do try:%d",
			// 				bindPodID, k8sResourceID, k8sAddrBindInfo.IP, err.Error(), replicas)
			// 		} else {
			// 			updateRows, _ := updateRes.RowsAffected()
			// 			// if err != nil {
			// 			// 	calm_utils.Errorf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s RowsAffected failed. err:%s. do try:%d",
			// 			// 		bindPodID, k8sResourceID, k8sAddrBindInfo.IP, err.Error(), replicas)
			// 			// 	continue
			// 			// }

			// 			calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_podid=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s successed. updateRows:[%d]",
			// 				bindPodID, k8sResourceID, k8sAddrBindInfo.IP, updateRows)

			// 			k8sAddrInfo = new(proto.K8SAddrInfo)
			// 			k8sAddrInfo.IP = k8sAddrBindInfo.IP
			// 			k8sAddrInfo.MacAddr = k8sAddrBindInfo.MacAddr
			// 			k8sAddrInfo.NetRegionalID = k8sAddrBindInfo.NetRegionalID
			// 			k8sAddrInfo.SubNetID = k8sAddrBindInfo.SubNetID
			// 			k8sAddrInfo.SubNetGatewayAddr = k8sAddrBindInfo.SubNetGatewayAddr
			// 			break
			// 		}
			// 	}
			// 	time.Sleep(time.Second)
			// 	replicas--
			// }
		} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
		}
		return nil
	})

	if k8sAddrInfo != nil {
		calm_utils.Infof("k8sResourceID:[%s] k8sResourceType:[%s] bindPodID:[%s] bind Addr:[%s]", k8sResourceID,
			k8sResourceType.String(), bindPodID, litter.Sdump(k8sAddrInfo))
	} else {
		calm_utils.Errorf("k8sResourceID:[%s] k8sResourceType:[%s] bindPodID:[%s] bind Addr failed", k8sResourceID,
			k8sResourceType.String(), bindPodID)
	}

	return k8sAddrInfo, err
}

// UnbindAddrInfoWithK8SPodID 地址和k8s资源解绑
func (msm *mysqlStoreMgr) UnbindAddrInfoWithK8SPodID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	unBindPodID string) error {

	return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=0 WHERE k8sresource_id=? AND bind_podid=?",
			k8sResourceID, unBindPodID)
		if err != nil {
			err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE k8sresource_id=%s AND bind_podid=%s failed.",
				k8sResourceID, unBindPodID)
			calm_utils.Error(err.Error())
			return err
		}
		updateRows, _ := updateRes.RowsAffected()
		calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE k8sresource_id=%s AND bind_podid=%s successed. updateRows:%d.",
			k8sResourceID, unBindPodID, updateRows)

		// 释放该IP
		var recycleIP, portID string
		ipBindRow := msm.dbMgr.QueryRow("SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_podid=?",
			k8sResourceID, unBindPodID)
		err = ipBindRow.Scan(&recycleIP, &portID)
		if err != nil {
			err = errors.Wrapf(err, "SELECT ip, port_id FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s AND bind_podid=%s failed.",
				k8sResourceID, unBindPodID)
			calm_utils.Error(err.Error())
			return err
		}

		// 判断该条记录是否要回收
		if k8sResourceType == proto.K8SApiResourceKindDeployment {
			var scaleDownMark table.TblK8SScaleDownMarkS
			err = msm.dbMgr.Get(&scaleDownMark, "SELECT * FROM tbl_K8SScaleDownMark WHERE k8sresource_id=? LIMIT 1", k8sResourceID)
			if err != nil {
				if !strings.Contains(err.Error(), "no rows in result set") {
					err = errors.Wrapf(err, "SELECT * FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s LIMIT 1", k8sResourceID)
					calm_utils.Error(err.Error())
					return err
				}
				// 没有标记，无需处理
				return nil
			}

			// 删除标记
			delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=? AND recycle_mark_id=?",
				k8sResourceID, scaleDownMark.RecycleMarkID)
			if err != nil {
				err = errors.Wrapf(err, "DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s AND recycle_mark_id=%s failed",
					k8sResourceID, scaleDownMark.RecycleMarkID)
				calm_utils.Error(err.Error())
				return err
			}

			delRows, _ := delRes.RowsAffected()
			if delRows != 1 {
				err = errors.Wrapf(err, "DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=%s AND recycle_mark_id=%s is incorrect, delRows:%d",
					k8sResourceID, scaleDownMark.RecycleMarkID, delRows)
				calm_utils.Error(err)
				return err
			}

			// 释放该条记录
			msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_podid=?",
				k8sResourceID, unBindPodID)
			// NSP回收
			nsp.NSPMgr.ReleaseAddrResources(portID)
		} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
			// TODO:
		}
		return nil
	})
}

// CheckRecycledResources 检查对应资源是否存在，bool = true存在，int=副本数量
func (msm *mysqlStoreMgr) CheckRecycledResources(k8sResourceID string) (bool, int, error) {
	// 在tbl_K8SResourceIPRecycle查询
	var replicas int
	err := msm.dbMgr.Get(&replicas, "SELECT replicas FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=?", k8sResourceID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// 地址资源不在待回收表中，
			calm_utils.Infof("%s not in tbl_K8SResourceIPRecycle, must be call NSP get addr resource.", k8sResourceID)
			return false, 0, nil
		}
		return false, 0, errors.Wrapf(err, "CheckRecycledResources %s failed.", k8sResourceID)
	}

	// 说明存在待回收表，删除该记录，恢复
	delRes, err := msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPRecycle WHERE k8sresource_id=?", k8sResourceID)
	if err != nil {
		err = errors.Wrapf(err, "DELETE tbl_K8SResourceIPRecycle by %s failed", k8sResourceID)
		calm_utils.Error(err.Error())
		return false, 0, err
	}

	delRows, _ := delRes.RowsAffected()
	if delRows != 1 {
		err = errors.Wrapf(err, "DELETE tbl_K8SResourceIPRecycle by %s delRows:%d is incorrect", k8sResourceID, delRows)
		calm_utils.Error(err)
		return false, 0, err
	} else {
		calm_utils.Infof("Restored resources:%s to be recovered", k8sResourceID)
	}

	return true, replicas, nil
}

// AddK8SResourceAddressToRecycle 加入回收站，待租期到期回收
func (msm *mysqlStoreMgr) AddK8SResourceAddressToRecycle(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType) error {
	// 查询已经分配的地址数量
	var k8sResourceReplicas int
	err := msm.dbMgr.Get(&k8sResourceReplicas, "SELECT COUNT(*) FROM tbl_K8SResourceIPBind WHERE k8sresource_id=?", k8sResourceID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// 地址资源不在绑定对应关系表中
			err = errors.Wrapf(err, "%s no leated records in tbl_K8SResourceIPBind.", k8sResourceID)
			calm_utils.Error(err.Error())
			return err
		}
		err = errors.Wrapf(err, "SELECT COUNT(*) FROM tbl_K8SResourceIPBind by k8sresource_id:%s failed.", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}

	calm_utils.Debugf("%s have %d count records in tbl_K8SResourceIPBind", k8sResourceID, k8sResourceReplicas)

	if k8sResourceReplicas == 0 {
		err = errors.Errorf("Empty Addr record for k8sresource_id:[%s]", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}

	// 插入租期回收表中
	k8sResourceIPRecycleRecord := new(table.TblK8SResourceIPRecycleS)
	k8sResourceIPRecycleRecord.SrvInstanceName = msm.opts.SrvInstID
	k8sResourceIPRecycleRecord.K8SResourceID = k8sResourceID
	k8sResourceIPRecycleRecord.K8SResourceType = int(k8sResourceType)
	k8sResourceIPRecycleRecord.Replicas = k8sResourceReplicas
	k8sResourceIPRecycleRecord.CreateTime = time.Now()
	k8sResourceIPRecycleRecord.NSPResourceReleaseTime = k8sResourceIPRecycleRecord.CreateTime.Add(config.GetK8SResourceAddrLeasePeriodSecs() * time.Second)
	k8sResourceIPRecycleRecord.RecycleObjectID = uuid.New().String()

	_, err = msm.dbMgr.Exec(`INSERT INTO tbl_K8SResourceIPRecycle 
	(srv_instance_name, k8sresource_id, k8sresource_type, replicas, create_time, nspresource_release_time, recycle_object_id) VALUES 
	(?, ?, ?, ?, ?, ?, ?)`,
		k8sResourceIPRecycleRecord.SrvInstanceName,
		k8sResourceIPRecycleRecord.K8SResourceID,
		k8sResourceIPRecycleRecord.K8SResourceType,
		k8sResourceIPRecycleRecord.Replicas,
		k8sResourceIPRecycleRecord.CreateTime,
		k8sResourceIPRecycleRecord.NSPResourceReleaseTime,
		k8sResourceIPRecycleRecord.RecycleObjectID,
	)
	if err != nil {
		err = errors.Wrapf(err, "INSERT tbl_K8SResourceIPRecycle k8sResourceID:%s failed.", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}

	// 加入定时堆管理
	msm.addrResourceLeasePeriodMgr.AddLeaseRecyclingRecord(k8sResourceIPRecycleRecord)

	calm_utils.Debugf("INSERT tbl_K8SResourceIPRecycle k8sResourceID:%s successed.", k8sResourceID)
	return nil
}

func (msm *mysqlStoreMgr) ScaleDownK8SResourceAddrs(k8sResourceID string, scaleDownSize int) error {

	calm_utils.Debugf("k8sResourceID:%s scaleDownSize:%d", k8sResourceID, scaleDownSize)
	return nil
}

func (msm *mysqlStoreMgr) AddScaleDownMarked(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	originalReplicas int, scaleDownSize int) error {

	if k8sResourceType == proto.K8SApiResourceKindDeployment {
		return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
			now := time.Now()
			tx := msm.dbMgr.MustBegin()
			// tbl_K8SScaleDownMark 插入
			for scaleDownSize > 0 {
				recycleMarkID := uuid.New().String()
				_, err := tx.Exec(`INSERT INTO tbl_K8SScaleDownMark (recycle_mark_id, k8sresource_id, k8sresource_type, create_time) VALUES (?, ?, ?, ?)`,
					recycleMarkID,
					k8sResourceID,
					k8sResourceType,
					now,
				)
				if err != nil {
					err = errors.Wrapf(err, "INSERT tbl_K8SScaleDownMark recycleMarkID:%s k8sResourceID:%s %d failed.",
						recycleMarkID, k8sResourceID, scaleDownSize)
					calm_utils.Error(err.Error())
					tx.Rollback()
					return err
				} else {
					calm_utils.Debugf("INSERT tbl_K8SScaleDownMark VALUES(%s, %s, %s, %s) successed.",
						recycleMarkID, k8sResourceID, k8sResourceType.String(), now.String())
				}
				scaleDownSize--
			}
			tx.Commit()
			return nil
		})
	}
	return nil
}
