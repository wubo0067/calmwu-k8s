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
	"fmt"
	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/config"
	"pci-ipresmgr/pkg/ipresmgr/k8s"
	"pci-ipresmgr/pkg/ipresmgr/nsp"
	"pci-ipresmgr/table"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
func (msm *mysqlStoreMgr) SetAddrInfosToK8SResourceID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	k8sAddrInfos []*proto.K8SAddrInfo, offset int) error {

	if k8sResourceType == proto.K8SApiResourceKindDeployment {
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
	} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
		return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
			tx := msm.dbMgr.MustBegin()

			for index, k8sAddrInfo := range k8sAddrInfos {
				calm_utils.Debugf("%d k8sResourceID[%s] k8sResourceType[%s] k8sAddrInfo:%+v", index, k8sResourceID,
					k8sResourceType.String(), k8sAddrInfo)

				insRes := tx.MustExec(`INSERT INTO tbl_K8SResourceIPBind
			(k8sresource_id, k8sresource_type, ip, mac, netregional_id, subnet_id, port_id, subnetgatewayaddr, alloc_time, is_bind, bind_poduniquename) VALUES
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
					fmt.Sprintf("%s-%d", k8sResourceID, index+offset), // 直接分配podid，因为statefulset的pod是固定的
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
	err := errors.Errorf("k8sResourceID[%s] k8sResourceType[%d] is invalid!", k8sResourceID, k8sResourceType)
	calm_utils.Error(err)
	return err
}

// BindAddrInfoWithK8SPodUniqueName 获取一个地址信息，和k8s资源绑定
func (msm *mysqlStoreMgr) BindAddrInfoWithK8SPodUniqueName(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	podUniqueName string) (*proto.K8SAddrInfo, error) {

	calm_utils.Debugf("k8sResourceID:[%s] k8sResourceType:[%s] podUniqueName:[%s]", k8sResourceID, k8sResourceType.String(),
		podUniqueName)

	var k8sAddrInfo *proto.K8SAddrInfo

	err := msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		if k8sResourceType == proto.K8SApiResourceKindDeployment ||
			k8sResourceType == proto.K8SApiResourceKindStatefulSet {

			// 判断这个pod是否已经绑定了，statefulset和deployment是不同的
			// statefulset的podname可以是相同的，如果重复绑定肯定是报错，要返回失败
			var checkExistk8sAddrBindInfo table.TblK8SResourceIPBindS
			err := msm.dbMgr.Get(&checkExistk8sAddrBindInfo,
				"SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND k8sresource_type=? AND is_bind=1 AND bind_poduniquename=? LIMIT 1",
				k8sResourceID, k8sResourceType, podUniqueName)
			if err == nil {
				if k8sResourceType == proto.K8SApiResourceKindDeployment {
					// 该pod已经绑定，直接返回
					k8sAddrInfo = new(proto.K8SAddrInfo)
					k8sAddrInfo.IP = checkExistk8sAddrBindInfo.IP
					k8sAddrInfo.MacAddr = checkExistk8sAddrBindInfo.MacAddr
					k8sAddrInfo.NetRegionalID = checkExistk8sAddrBindInfo.NetRegionalID
					k8sAddrInfo.SubNetID = checkExistk8sAddrBindInfo.SubNetID
					k8sAddrInfo.SubNetGatewayAddr = checkExistk8sAddrBindInfo.SubNetGatewayAddr
					k8sAddrInfo.PortID = checkExistk8sAddrBindInfo.PortID
					calm_utils.Warnf("k8sResourceType:[%s] k8sResourceID:[%s] bindPod:[%s] is already occupied address:%s resources", k8sResourceID, podUniqueName, checkExistk8sAddrBindInfo.IP)
					return nil
				} else {
					// statefulset不允许重复绑定
					err := errors.Errorf("k8sResourceType:[%s] k8sResourceID:[%s] bindPod:[%s] is already occupied address:%s resources", k8sResourceID, podUniqueName, checkExistk8sAddrBindInfo.IP)
					calm_utils.Error(err.Error())
					return err
				}
			}

			// https://www.cnblogs.com/diegodu/p/9239200.html 用串行化事务，gap lock
			tx, err := msm.dbMgr.BeginTxx(context.Background(), &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
			})
			//tx, err := msm.dbMgr.Begin()
			if err != nil {
				err := errors.Wrapf(err, "k8sResourceID:%s bindPod:%s begin transaction failed.", k8sResourceID, podUniqueName)
				calm_utils.Error(err.Error())
				return err
			}

			var transactionFlag int
			defer func(flag *int) {
				if *flag == 0 {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Commit", k8sResourceID, podUniqueName)
					tx.Commit()
				} else {
					calm_utils.Debugf("k8sResourceID:%s bindPod:%s tbl_K8SResourceIPBind SELECT FOR UPDATE Rollback", k8sResourceID, podUniqueName)
					tx.Rollback()
				}
			}(&transactionFlag)

			queryBindIPSql := ""
			if k8sResourceType == proto.K8SApiResourceKindDeployment {
				queryBindIPSql = fmt.Sprintf("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id='%s AND k8sresource_type=%d FOR UPDATE",
					k8sResourceID, int(k8sResourceType))
			} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
				queryBindIPSql = fmt.Sprintf("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id='%s AND k8sresource_type=%d AND bind_poduniquename='%s' FOR UPDATE",
					k8sResourceID, int(k8sResourceType), podUniqueName)
			}

			selRows, err := tx.Queryx(queryBindIPSql)
			if err != nil {
				err = errors.Wrapf(err, "%s failed", queryBindIPSql)
				calm_utils.Error(err.Error())
				transactionFlag = -1
				return err
			}

			var isFindUnbindAddr bool = false
			var k8sAddrBindInfo table.TblK8SResourceIPBindS
			// 循环过滤出未绑定的地址资源
			for selRows.Next() {
				err := selRows.StructScan(&k8sAddrBindInfo)
				if err != nil {
					err = errors.Wrapf(err, "%s StructScan failed.", queryBindIPSql)
					calm_utils.Error(err.Error())
					transactionFlag = -1
					selRows.Close()
					return err
				}
				if k8sAddrBindInfo.IsBind == 0 {
					isFindUnbindAddr = true
					break
				}
			}
			selRows.Close()

			if !isFindUnbindAddr {
				err = errors.Errorf("%s not found unBindAddrInfo.",
					k8sResourceID, k8sResourceType)
				calm_utils.Error(err.Error())
				transactionFlag = -1
				return err
			}

			calm_utils.Debugf("k8sResourceID:%s bindPod:%s k8sAddrBindInfo:%s", k8sResourceID, podUniqueName, litter.Sdump(&k8sAddrBindInfo))

			currTime := time.Now()
			updateRes, err := tx.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=?, bind_time=? WHERE k8sresource_id=? AND port_id=?",
				podUniqueName, currTime, k8sResourceID, k8sAddrBindInfo.PortID)
			if err != nil {
				err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET is_bind=1, bind_poduniquename=%s, bind_time=%s WHERE k8sresource_id=%s AND port_id=%s. tx Exec UPDATE failed.",
					podUniqueName, currTime.String(), k8sResourceID, k8sAddrBindInfo.PortID)
				calm_utils.Error(err.Error())
				transactionFlag = -1
				return err
			}

			updateRowCount, _ := updateRes.RowsAffected()

			calm_utils.Debugf("k8sResourceID:%s bindPod:%s updateRowCount:%d\n", k8sResourceID, podUniqueName, updateRowCount)

			k8sAddrInfo = new(proto.K8SAddrInfo)
			k8sAddrInfo.IP = k8sAddrBindInfo.IP
			k8sAddrInfo.MacAddr = k8sAddrBindInfo.MacAddr
			k8sAddrInfo.NetRegionalID = k8sAddrBindInfo.NetRegionalID
			k8sAddrInfo.SubNetID = k8sAddrBindInfo.SubNetID
			k8sAddrInfo.SubNetGatewayAddr = k8sAddrBindInfo.SubNetGatewayAddr
			k8sAddrInfo.PortID = k8sAddrBindInfo.PortID
		}
		return nil
	})

	if k8sAddrInfo != nil {
		calm_utils.Infof("k8sResourceID:[%s] k8sResourceType:[%s] podUniqueName:[%s] bind Addr:[%s]", k8sResourceID,
			k8sResourceType.String(), podUniqueName, litter.Sdump(k8sAddrInfo))
	} else {
		// TODO 发送告警
		calm_utils.Errorf("k8sResourceID:[%s] k8sResourceType:[%s] podUniqueName:[%s] bind Addr failed", k8sResourceID,
			k8sResourceType.String(), podUniqueName)
	}

	return k8sAddrInfo, err
}

// UnbindAddrInfoWithK8SPodID 地址和k8s资源解绑
func (msm *mysqlStoreMgr) UnbindAddrInfoWithK8SPodID(k8sResourceType proto.K8SApiResourceKindType, podUniqueName string) error {

	// podUniqueName 是唯一索引
	return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
		updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SResourceIPBind SET is_bind=0 WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
		if err != nil {
			err = errors.Wrapf(err, "UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE bind_poduniquename=%s LIMIT 1 failed.", podUniqueName)
			calm_utils.Error(err.Error())
			return err
		}
		updateRows, _ := updateRes.RowsAffected()
		calm_utils.Debugf("UPDATE tbl_K8SResourceIPBind SET bind=0 WHERE bind_poduniquename=%s successed. updateRows:%d.", podUniqueName, updateRows)

		// 释放该IP
		var k8sResourceID, recycleIP, portID string
		ipBindRow := msm.dbMgr.QueryRow("SELECT k8sresource_id, ip, port_id FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
		err = ipBindRow.Scan(&k8sResourceID, &recycleIP, &portID)
		if err != nil {
			err = errors.Wrapf(err, "SELECT k8sresource_id, ip, port_id FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=%s LIMIT 1 failed.", podUniqueName)
			calm_utils.Error(err.Error())
			return err
		}

		// 判断该条记录是否要回收
		if k8sResourceType == proto.K8SApiResourceKindDeployment {
			// 如果对deployment做了scaledown，需要判断释放的条数，并释放该条记录
			updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SScaleDownMark SET scaledown_count = scaledown_count-1 WHERE k8sresource_id=? AND scaledown_count > 0 LIMIT 1",
				k8sResourceID)
			if err != nil {
				// 没有scaledown记录，更新不成功，不用回收
				calm_utils.Infof("BindPodID:%s not set scaledown flag, No immediate release required.", podUniqueName)
				return nil
			}
			updateRowCount, _ := updateRes.RowsAffected()
			if updateRowCount != 1 {
				calm_utils.Infof("BindPodID:%s not set scaledown flag, No immediate release required.", podUniqueName)
				// 有记录，但是没有更新，说明scaledown已经结束，是个正常的删除pod行为，需要将该条记录删除掉
				msm.dbMgr.Exec("DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=?", k8sResourceID)
				return nil
			}
			// NSP回收
			calm_utils.Debugf("Deployment BindPodID:%s set scaledown flag, so release immediately", podUniqueName)
			msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
			nsp.NSPMgr.ReleaseAddrResources(portID)
		} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
			// statefulset要精确计算，该回收那些，避免在scaledown的时候做了delete操作，这样会导致释放了不该的ip地址
			podNameSeq, _ := strconv.Atoi(podUniqueName[strings.LastIndexByte(podUniqueName, '-')+1:])
			// 必须保证释放的podseq是大于等于当前副本数的
			updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SScaleDownMark SET scaledown_count = scaledown_count-1 WHERE k8sresource_id=? AND scaledown_count > 0 AND current_replicas <= ? LIMIT 1",
				k8sResourceID, podNameSeq)
			if err != nil {
				// 没有scaledown记录，更新不成功，不用回收
				calm_utils.Infof("BindPodID:%s not set scaledown flag, No immediate release required.", podUniqueName)
				return nil
			}
			updateRowCount, _ := updateRes.RowsAffected()
			if updateRowCount != 1 {
				calm_utils.Infof("BindPodID:%s not set scaledown flag, No immediate release required.", podUniqueName)
				// 有记录，但是没有更新，说明scaledown已经结束，是个正常的删除pod行为，需要将该条记录删除掉
				msm.dbMgr.Exec("DELETE FROM tbl_K8SScaleDownMark WHERE k8sresource_id=?", k8sResourceID)
				return nil
			}
			// NSP回收
			calm_utils.Debugf("Statefulset BindPodID:%s set scaledown flag, so release immediately", podUniqueName)
			msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
			nsp.NSPMgr.ReleaseAddrResources(portID)
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
// TODO: statefulset的释放
func (msm *mysqlStoreMgr) ReduceK8SResourceAddrs(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, reduceCount int) error {
	// 找出所有对应地址，见解绑中的地址进行回收，如果数量不够，就等待，等待超时就失败
	reduceRows, err := msm.dbMgr.Queryx("SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=?", k8sResourceID)
	if err != nil {
		errors.Wrapf(err, "SELECT * FROM tbl_K8SResourceIPBind WHERE k8sresource_id=%s failed.", k8sResourceID)
		calm_utils.Error(err.Error())
		return err
	}

	reduceK8SBindAddrs := make([]*table.TblK8SResourceIPBindS, 0)
	unBindCount := 0
	for reduceRows.Next() {
		k8sBindAddr := new(table.TblK8SResourceIPBindS)
		err = reduceRows.StructScan(k8sBindAddr)
		if err != nil {
			errors.Wrapf(err, "StructScan for tbl_K8SResourceIPBind record failed.")
			calm_utils.Error(err.Error())
			return err
		}
		reduceK8SBindAddrs = append(reduceK8SBindAddrs, k8sBindAddr)
		if k8sBindAddr.IsBind == 0 {
			unBindCount++
		}
	}

	redceK8SBindAddrsLen := len(reduceK8SBindAddrs)

	calm_utils.Infof("k8sResourceID:%s reduceCount:%d UnBind Addr count:%d total count:%d", k8sResourceID, reduceCount, unBindCount, redceK8SBindAddrsLen)

	if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
		// 排序，序号大的在前
		sort.Slice(reduceK8SBindAddrs, func(i, j int) bool {
			return reduceK8SBindAddrs[i].BindPodUniqueName.String > reduceK8SBindAddrs[j].BindPodUniqueName.String
		})

		for index := range reduceK8SBindAddrs {
			calm_utils.Debugf("%d %s", index, litter.Sdump(reduceK8SBindAddrs[index]))
		}
	}

	if unBindCount < reduceCount {
		// TODO: 告警，已经解绑的pod数量小于缩容的数量
		// 去查询还有哪些没有释放的pod状态，node状态
		for _, k8sBindAddr := range reduceK8SBindAddrs {
			k8s.DefaultK8SClient.GetPodAndNodeStatus(k8sResourceID, k8sBindAddr.PortID)
		}
		err = errors.Errorf("k8sResourceID:%s Failure to reduce the number[%d] of IPs", k8sResourceID, reduceCount)
		calm_utils.Error(err.Error())
		return err
	} else {
		if k8sResourceType == proto.K8SApiResourceKindDeployment {
			// 找到reduce count的unbind地址进行释放
			for _, k8sBindAddr := range reduceK8SBindAddrs {
				if k8sBindAddr.IsBind == 0 {
					calm_utils.Infof("Deployment POD k8sResourceID:%s BindPodUniqueName:%s ip:%s portID:%s address will be recycled and returned to nsp",
						k8sResourceID, k8sBindAddr.BindPodUniqueName, k8sBindAddr.IP, k8sBindAddr.PortID)
					// 删除该条记录
					msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_poduniquename=? LIMIT 1",
						k8sResourceID, k8sBindAddr.BindPodUniqueName)
					// NSP回收
					nsp.NSPMgr.ReleaseAddrResources(k8sBindAddr.PortID)
					reduceCount--
					if reduceCount == 0 {
						break
					}
				}
			}
		} else {
			// statefulset，按排序释放
			for index := 0; index < reduceCount; index++ {
				k8sBindAddr := reduceK8SBindAddrs[index]
				if k8sBindAddr.IsBind != 0 {
					// TODO 告警
				}
				calm_utils.Infof("StatefulSet POD k8sResourceID:%s BindPodUniqueName:%s ip:%s portID:%s address will be recycled and returned to nsp",
					k8sResourceID, k8sBindAddr.BindPodUniqueName, k8sBindAddr.IP, k8sBindAddr.PortID)
				// 删除该条记录
				msm.dbMgr.Exec("DELETE FROM tbl_K8SResourceIPBind WHERE k8sresource_id=? AND bind_poduniquename=? LIMIT 1",
					k8sResourceID, k8sBindAddr.BindPodUniqueName)
				// NSP回收
				nsp.NSPMgr.ReleaseAddrResources(k8sBindAddr.PortID)
			}
		}

	}
	return nil
}

func (msm *mysqlStoreMgr) AddScaleDownMarked(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType,
	currReplicas int, scaleDownSize int) error {

	if k8sResourceType == proto.K8SApiResourceKindDeployment ||
		k8sResourceType == proto.K8SApiResourceKindStatefulSet {
		return msm.dbSafeExec(context.Background(), func(ctx context.Context) error {
			calm_utils.Debugf("k8sResourceID:%s k8sResourceType:%s scaleDownSize:%d flag",
				k8sResourceID, k8sResourceType.String(), scaleDownSize)

			createTime := time.Now()
			insertRes, err := msm.dbMgr.Exec("INSERT INTO tbl_K8SScaleDownMark (k8sresource_id, k8sresource_type, current_replicas, scaledown_count, create_time) VALUES (?, ?, ?, ?)",
				k8sResourceID, int(k8sResourceType), currReplicas, scaleDownSize, createTime)
			if err != nil {
				if strings.Contains(err.Error(), "for key 'PRIMARY'") {
					// 直接更新数量
					calm_utils.Infof("k8sResourceID:%s already in tbl_K8SScaleDownMark, so add scaledown_count:%d", k8sResourceID, scaleDownSize)
					updateRes, err := msm.dbMgr.Exec("UPDATE tbl_K8SScaleDownMark SET scaledown_count=scaledown_count+?, current_replicas=? WHERE k8sresource_id=?", scaleDownSize, currReplicas, k8sResourceID)
					if err != nil {
						err = errors.Wrapf(err, "UPDATE tbl_K8SScaleDownMark SET scaledown_count=scaledown_count+%d, current_replicas=%d WHERE k8sresource_id=%s failed.",
							scaleDownSize, currReplicas, k8sResourceID)
						calm_utils.Error(err.Error())
						return err
					}
					updateRows, _ := updateRes.RowsAffected()
					if updateRows != 1 {
						err = errors.Wrapf(err, "UPDATE tbl_K8SScaleDownMark SET scaledown_count=scaledown_count+%d, current_replicas=%d WHERE k8sresource_id=%s no affect, updateRows:%d.",
							scaleDownSize, currReplicas, k8sResourceID, updateRows)
						calm_utils.Error(err.Error())
						return err
					}
					calm_utils.Debugf("UPDATE tbl_K8SScaleDownMark SET scaledown_count=scaledown_count+%d, current_replicas=%d WHERE k8sresource_id=%s successed.",
						scaleDownSize, currReplicas, k8sResourceID)
					return nil
				}
				err = errors.Wrapf(err, "INSERT INTO tbl_K8SScaleDownMark (k8sresource_id, k8sresource_type, current_replicas, scaledown_count, create_time) VALUES (%s, %s, %d, %d, %s) failed.",
					k8sResourceID, k8sResourceType.String(), currReplicas, scaleDownSize, createTime.String())
				calm_utils.Error(err.Error())
				return err
			}
			insertRows, _ := insertRes.RowsAffected()
			if insertRows != 1 {
				err = errors.Errorf("INSERT INTO tbl_K8SScaleDownMark (k8sresource_id, k8sresource_type, current_replicas, scaledown_count, create_time) VALUES (%s, %s, %d, %d, %s) insertRows[%d] Not equal to 1.",
					k8sResourceID, k8sResourceType.String(), currReplicas, scaleDownSize, createTime.String(), insertRows)
				calm_utils.Error(err.Error())
				return err
			}
			calm_utils.Debugf("INSERT INTO tbl_K8SScaleDownMark (k8sresource_id, k8sresource_type, current_replicas, scaledown_count, create_time) VALUES (%s, %s, %d, %d, %s) insertRows[%d] successed.",
				k8sResourceID, k8sResourceType.String(), currReplicas, scaleDownSize, createTime.String(), insertRows)
			return nil
		})
	}
	return nil
}

func (msm *mysqlStoreMgr) QueryK8SResourceKindByPodUniqueName(podUniqueName string) proto.K8SApiResourceKindType {
	var k8sResourceType int
	err := msm.dbMgr.Get(&k8sResourceType, "SELECT k8sresource_type FROM tbl_K8SResourceIPBind WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
	if err == nil {
		kindType := proto.K8SApiResourceKindType(k8sResourceType)
		calm_utils.Debugf("podUniqueName:%s in tbl_K8SResourceIPBind, type is %s", podUniqueName, kindType.String())
		return kindType
	}

	err = msm.dbMgr.Get(&k8sResourceType, "SELECT k8sresource_type FROM tbl_K8SJobIPBind WHERE bind_poduniquename=? LIMIT 1", podUniqueName)
	if err == nil {
		kindType := proto.K8SApiResourceKindType(k8sResourceType)
		calm_utils.Debugf("podUniqueName:%s in tbl_K8SResourceIPBind, type is %s", podUniqueName, kindType.String())
		return kindType
	}

	return proto.K8SApiResourceKindUnknown
}
