/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:59:46
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-01 10:21:19
 */

// 地址租期管理

package mysql

import (
	"context"
	"pci-ipresmgr/table"
	"sync"

	"github.com/jmoiron/sqlx"
)

var _ AddrResourceLeasePeriodMgrItf = &mysqlAddrResourceLeasePeriodMgr{}

// AddrResourceLeasePeriodMgrItf 回收接口
type AddrResourceLeasePeriodMgrItf interface {
	// 启动
	Start() error
	// 停止
	Stop()
	// 添加租期回收对象
	AddLeaseRecyclingRecord(record *table.TblK8SResourceIPRecycleS)
	// 删除租期回收对象
	DelLeaseRecyclingRecord(k8sResourceID string) *table.TblK8SResourceIPRecycleS
}

// implement mysqlAddrResourceLeasePeriodMgr
type mysqlAddrResourceLeasePeriodMgr struct {
	ctx   context.Context
	dbMgr *sqlx.DB
	guard sync.Mutex
}

func (malm *mysqlAddrResourceLeasePeriodMgr) Start() error {
	return nil
}

func (malm *mysqlAddrResourceLeasePeriodMgr) Stop() {
	return
}

func (malm *mysqlAddrResourceLeasePeriodMgr) AddLeaseRecyclingRecord(record *table.TblK8SResourceIPRecycleS) {
	return
}

func (malm *mysqlAddrResourceLeasePeriodMgr) DelLeaseRecyclingRecord(k8sResourceID string) *table.TblK8SResourceIPRecycleS {
	return nil
}

// NewAddrResourceLeasePeriodMgr 构造地址租期管理对象
func NewAddrResourceLeasePeriodMgr(ctx context.Context, dbMgr *sqlx.DB) AddrResourceLeasePeriodMgrItf {
	addrResourceLeasePeriodMgr := new(mysqlAddrResourceLeasePeriodMgr)
	addrResourceLeasePeriodMgr.ctx = ctx
	addrResourceLeasePeriodMgr.dbMgr = dbMgr
	return addrResourceLeasePeriodMgr
}
