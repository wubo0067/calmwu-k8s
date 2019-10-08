/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:52:15
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:53:59
 */

package mysql

import (
	"context"
	"database/sql"
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/k8s"
	"pci-ipresmgr/pkg/ipresmgr/nsp"
	"pci-ipresmgr/table"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

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

// BindAddrInfoWithK8SPodUniqueName 获取一个地址信息，和k8s资源绑定
func (msm *mysqlStoreMgr) BindAddrInfoWithK8SPodUniqueName(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	bindPodUniqueName string) (*proto.K8SAddrInfo, error) {

	var k8sAddrInfo *proto.K8SAddrInfo

	err := msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		if k8sResourceType == proto.K8SApiResourceKindDeployment {

			// 判断这个pod是否已经绑定了
			var checkExistk8sAddrBindInfo table.TblK8SResourceIPBindS
			err := msm.dbMgr.Get(&checkExistk8sAddrBindInfo,
				"SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=1 AND bind_poduniquename=? LIMIT 1",
				k8sResourceID, k8sResourceType, bindPodUniqueName)
			if err == nil {
				// 该pod已经绑定，直接返回
				k8sAddrInfo = new(proto.K8SAddrInfo)
				k8sAddrInfo.IP = checkExistk8sAddrBindInfo.IP
				k8sAddrInfo.MacAddr = checkExistk8sAddrBindInfo.MacAddr
				k8sAddrInfo.NetRegionalID = checkExistk8sAddrBindInfo.NetRegionalID
				k8sAddrInfo.SubNetID = checkExistk8sAddrBindInfo.SubNetID
				k8sAddrInfo.SubNetGatewayAddr = checkExistk8sAddrBindInfo.SubNetGatewayAddr
				k8sAddrInfo.PortID = checkExistk8sAddrBindInfo.PortID
				calm_utils.Warnf("k8sResourceID:%s bindPod:%s is already occupied address:%s resources", k8sResourceID, bindPodUniqueName, checkExistk8sAddrBindInfo.IP)
				return nil
			}

			// https://www.cnblogs.com/diegodu/p/9239200.html 用串行化事务，gap lock
			tx, err := msm.dbMgr.BeginTx(context.Background(), &sql.TxOptions{
				Isolation: sql.LevelSerializable,
			})
			if err != nil {
				err := errors.Wrapf(err, "k8sResourceID:%s bindPod:%s begin transaction failed.", k8sResourceID, bindPodUniqueName)
				calm_utils.Error(err.Error())
				return err
			}

			var transactionFlag int
			defer func(flag *int) {
				if *flag == 0 {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Commit", k8sResourceID, bindPodUniqueName)
					tx.Commit()
				} else {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Rollback", k8sResourceID, bindPodUniqueName)
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
				&k8sAddrBindInfo.SubNetGatewayAddr, &k8sAddrBindInfo.AllocTime, &k8sAddrBindInfo.IsBind, &k8sAddrBindInfo.BindPodUniqueName,
				&k8sAddrBindInfo.BindTime, &k8sAddrBindInfo.ScaledownFlag)

			if err != nil {
				// TODO 查询node，和所有pod的状态
				err = errors.Wrapf(err, "k8sResourceID:%s bindPod:%s Scan failed.", k8sResourceID, bindPodUniqueName)
				calm_utils.Error(err.Error())
				transactionFlag = -1
				return err
			}

			calm_utils.Debugf("k8sResourceID:%s bindPod:%s k8sAddrBindInfo:%s", k8sResourceID, bindPodUniqueName, litter.Sdump(&k8sAddrBindInfo))

			currTime := time.Now()
			updateRes, err := tx.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=?, bind_time=? WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 AND port_id=?",
				bindPodUniqueName, currTime, k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.PortID)
			if err != nil {
				err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=%s, bind_time=%s WHERE k8sresource_id=%s AND k8sresource_type=%d AND is_bind=0 AND port_id=%s. tx Exec UPDATE failed.",
					bindPodUniqueName, currTime.String(), k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.PortID)
				transactionFlag = -1
				return err
			}

			updateRowCount, _ := updateRes.RowsAffected()
			calm_utils.Debugf("k8sResourceID:%s bindPod:%s updateRowCount:%d\n", k8sResourceID, bindPodUniqueName, updateRowCount)

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
			// 		updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=?, bind_time=? WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=0 AND ip=?",
			// 			bindPodUniqueName, currTime, k8sResourceID, int(k8sResourceType), k8sAddrBindInfo.IP)
			// 		if err != nil {
			// 			calm_utils.Errorf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s failed. err:%s. do try:%d",
			// 				bindPodUniqueName, k8sResourceID, k8sAddrBindInfo.IP, err.Error(), replicas)
			// 		} else {
			// 			updateRows, _ := updateRes.RowsAffected()
			// 			// if err != nil {
			// 			// 	calm_utils.Errorf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s RowsAffected failed. err:%s. do try:%d",
			// 			// 		bindPodUniqueName, k8sResourceID, k8sAddrBindInfo.IP, err.Error(), replicas)
			// 			// 	continue
			// 			// }

			// 			calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=%s WHERE k8sresource_id=%s AND k8sresource_type=0 AND is_bind=0 AND ip=%s successed. updateRows:[%d]",
			// 				bindPodUniqueName, k8sResourceID, k8sAddrBindInfo.IP, updateRows)

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
		calm_utils.Infof("k8sResourceID:[%s] k8sResourceType:[%s] bindPodUniqueName:[%s] bind Addr:[%s]", k8sResourceID,
			k8sResourceType.String(), bindPodUniqueName, litter.Sdump(k8sAddrInfo))
	} else {
		// TODO 发送告警
		calm_utils.Errorf("k8sResourceID:[%s] k8sResourceType:[%s] bindPodUniqueName:[%s] bind Addr failed", k8sResourceID,
			k8sResourceType.String(), bindPodUniqueName)
	}

	return k8sAddrInfo, err
}

// UnbindAddrInfoWithK8SPodID 地址和k8s资源解绑
func (msm *mysqlStoreMgr) UnbindAddrInfoWithK8SPodID(k8sResourceType proto.K8SApiResourceKindType, unBindPodUniqueName string) error {

	return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=0 WHERE bind_poduniquename=? LIMIT 1", unBindPodUniqueName)
		if err != nil {
			err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE bind_poduniquename=%s LIMIT 1 failed.", unBindPodUniqueName)
			calm_utils.Error(err.Error())
			return err
		}
		updateRows, _ := updateRes.RowsAffected()
		calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE bind_poduniquename=%s successed. updateRows:%d.", unBindPodUniqueName, updateRows)

		// 释放该IP
		var recycleIP, portID string
		var scaleDownFlag int
		ipBindRow := msm.dbMgr.QueryRow("SELECT ip, port_id, scaledown_flag FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1", unBindPodUniqueName)
		err = ipBindRow.Scan(&recycleIP, &portID, &scaleDownFlag)
		if err != nil {
			err = errors.Wrapf(err, "SELECT ip, port_id, scaledown_flag FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=%s LIMIT 1 failed.", unBindPodUniqueName)
			calm_utils.Error(err.Error())
			return err
		}

		// 判断该条记录是否要回收
		if k8sResourceType == proto.K8SApiResourceKindDeployment {
			if scaleDownFlag == 1 {
				calm_utils.Infof("BindPodID:%s set scaledown flag, so release immediately", unBindPodUniqueName)
				// 释放该条记录
				msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=?", unBindPodUniqueName)
				// NSP回收
				nsp.NSPMgr.ReleaseAddrResources(portID)
			}
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

	calm_utils.Debugf("k8sResourceID:%s bind %d addresses in tbl_K8SResourceIPBind", k8sResourceID, k8sResourceReplicas)

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

// ReduceK8SResourceAddrs 租期恢复期间降低副本数量
func (msm *mysqlStoreMgr) ReduceK8SResourceAddrs(k8sResourceID string, reduceCount int) error {
	// 找出所有对应地址，见解绑中的地址进行回收，如果数量不够，就等待，等待超时就失败
	calm_utils.Debugf("k8sResourceID:%s reduceCount:%d", k8sResourceID, reduceCount)

	reduceRows, err := msm.dbMgr.Queryx("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=?", k8sResourceID)
	if err != nil {
		errors.Wrapf(err, "SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s failed.", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}

	reduceK8sBindAddrs := make([]*table.TblK8SResourceIPBindS, 0)
	unBindCount := 0
	for reduceRows.Next() {
		k8sBindAddr := new(table.TblK8SResourceIPBindS)
		err = reduceRows.StructScan(k8sBindAddr)
		if err != nil {
			errors.Wrapf(err, "StructScan for tbl_K8SResourceIPBind record failed.")
			calm_utils.Error(err.Error())
			return err
		}
		reduceK8sBindAddrs = append(reduceK8sBindAddrs, k8sBindAddr)
		if k8sBindAddr.IsBind == 0 {
			unBindCount++
		}
	}

	calm_utils.Infof("k8sResourceID:%s UnBind Addr count:%d", k8sResourceID, unBindCount)

	if unBindCount < reduceCount {
		// TODO: 告警
		// 去查询还有哪些没有释放的pod状态，node状态
		for _, k8sBindAddr := range reduceK8sBindAddrs {
			k8s.DefaultK8SClient.GetPodAndNodeStatus(k8sResourceID, k8sBindAddr.PortID)
		}
	} else {
		// 找到reduce count的unbind地址进行释放
		for _, k8sBindAddr := range reduceK8sBindAddrs {
			if k8sBindAddr.IsBind == 0 {
				calm_utils.Infof("k8sResourceID:%s BindPodUniqueName:%s ip:%s portID:%s address will be recycled and returned to nsp",
					k8sResourceID, k8sBindAddr.BindPodUniqueName, k8sBindAddr.IP, k8sBindAddr.PortID)
				// 删除该条记录
				msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_poduniquename=? LIMIT 1",
					k8sResourceID, k8sBindAddr.BindPodUniqueName)
				// NSP回收
				nsp.NSPMgr.ReleaseAddrResources(k8sBindAddr.PortID)
				unBindCount--
				if unBindCount == 0 {
					break
				}
			}
		}
	}
	return nil
}

func (msm *mysqlStoreMgr) AddScaleDownMarked(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	originalReplicas int, scaleDownSize int) error {

	if k8sResourceType == proto.K8SApiResourceKindDeployment {
		return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
			calm_utils.Debugf("Now Set k8sResourceID:%s k8sResourceType:%s scaleDownSize:%d flag",
				k8sResourceID, k8sResourceType.String(), scaleDownSize)

			updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET scaledown_flag=1 WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=1 LIMIT ?",
				k8sResourceID, int(k8sResourceType), scaleDownSize)
			if err != nil {
				err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET scaledown_flag=1 WHERE k8sresource_id=%s AND k8sresource_type=%d AND is_bind=1 LIMIT %d failed.",
					k8sResourceID, int(k8sResourceType), scaleDownSize)
				calm_utils.Error(err.Error())
				return err
			}
			updateRows, _ := updateRes.RowsAffected()
			calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET scaledown_flag=1 WHERE k8sresource_id=%s AND k8sresource_type=%d AND is_bind=1 LIMIT %d successed. updateRows:%d.",
				k8sResourceID, int(k8sResourceType), scaleDownSize, updateRows)
			return nil
		})
	}
	return nil
}

func (msm *mysqlStoreMgr) QueryK8SResourceKindByPodUniqueName(unBindPodUniqueName string) proto.K8SApiResourceKindType {
	var k8sResourceType int
	err := msm.dbMgr.Get(&k8sResourceType, `SELECT k8sresource_type FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1`, unBindPodUniqueName)
	if err != nil {
		calm_utils.Infof("unBindPodUniqueName:%s not in tbl_K8SResourceIPBind, so type is proto.K8SApiResourceKindLikeJob", unBindPodUniqueName)
		return proto.K8SApiResourceKindLikeJob
	}

	kindType := proto.K8SApiResourceKindType(k8sResourceType)
	calm_utils.Debugf("unBindPodUniqueName:%s in tbl_K8SResourceIPBind, type is %s", unBindPodUniqueName, kindType.String())

	return kindType
}
