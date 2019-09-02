/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:59:46
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-02 19:34:36
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

// addrLeasePeriodTimerNodeS 地址租期时间对象
type addrLeasePeriodTimerNodeS struct {
	record *table.TblK8SResourceIPRecycleS
	index  int
}

// addrLeasePeriodTimerHeap 地址租期最小堆
type addrLeasePeriodTimerHeap []*addrLeasePeriodTimerNodeS

func (ath addrLeasePeriodTimerHeap) Len() int {
	return len(ath)
}

func (ath addrLeasePeriodTimerHeap) Less(i, j int) bool {
	return ath[i].record.NSPResourceReleaseTime.Before(ath[j].record.NSPResourceReleaseTime)
}

func (ath addrLeasePeriodTimerHeap) Swap(i, j int) {
	ath[i], ath[j] = ath[j], ath[i]
	ath[i].index = i
	ath[j].index = j
}

func (ath *addrLeasePeriodTimerHeap) Push(timerNode interface{}) {
	count := len(*ath)
	addrLeasePeriodTimerNode := timerNode.(*addrLeasePeriodTimerNodeS)
	addrLeasePeriodTimerNode.index = count
	*ath = append(*ath, addrLeasePeriodTimerNode)
}

func (ath *addrLeasePeriodTimerHeap) Pop() interface{} {
	old := *ath
	n := len(old)
	addrLeasePeriodTimerNode := old[n-1]
	addrLeasePeriodTimerNode.index = -1 // for safety
	*ath = old[0 : n-1]
	return addrLeasePeriodTimerNode
}

//-----------------------------------------------------------------------------------

// implement mysqlAddrResourceLeasePeriodMgr
type mysqlAddrResourceLeasePeriodMgr struct {
	ctx       context.Context
	dbMgr     *sqlx.DB
	guard     sync.Mutex
	timerHeap *addrLeasePeriodTimerHeap
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
