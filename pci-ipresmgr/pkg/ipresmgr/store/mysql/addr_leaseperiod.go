/*
 * @Author: calm.wu
 * @Date: 2019-09-01 09:59:46
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-02 19:34:36
 */

// 地址租期管理

package mysql

import (
	"container/heap"
	"context"
	"pci-ipresmgr/table"
	"sync"
	"time"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
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

// 该方法没有使用，不会修改节点的时间
func (ath *addrLeasePeriodTimerHeap) update(node *addrLeasePeriodTimerNodeS, record *table.TblK8SResourceIPRecycleS) {
	node.record = record
	heap.Fix(ath, node.index)
}

//-----------------------------------------------------------------------------------

// implement mysqlAddrResourceLeasePeriodMgr
type mysqlAddrResourceLeasePeriodMgr struct {
	ctx       context.Context
	msm       *mysqlStoreMgr
	guard     sync.Mutex
	timerHeap addrLeasePeriodTimerHeap
}

func (malm *mysqlAddrResourceLeasePeriodMgr) Start() error {
	// heap初始化
	heap.Init(&malm.timerHeap)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
	L:
		for {
			select {
			case <-ticker.C:
				// 定时检测租期是否到期
				mostRecentTimeout, timeHeapSize := malm.getMostRecentTimeout()
				if mostRecentTimeout.IsZero() {
					calm_utils.Debug("mostRecentTimeout is zero, timerHeap is empty")
					continue
				}
				now := time.Now()
				calm_utils.Debugf("mostRecentTimeout:%s timeHeapSize:%d", mostRecentTimeout.String(), timeHeapSize)
				if now.After(mostRecentTimeout) {
					// 到期
					malm.timeCheckExpiration(now)
				}
			case <-malm.ctx.Done():
				calm_utils.Info("mysqlAddrResourceLeasePeriodMgr Ticker exit")
				break L
			}
		}
		return
	}()
	return nil
}

func (malm *mysqlAddrResourceLeasePeriodMgr) Stop() {
	calm_utils.Debug("mysqlAddrResourceLeasePeriodMgr stop.")
	return
}

func (malm *mysqlAddrResourceLeasePeriodMgr) AddLeaseRecyclingRecord(record *table.TblK8SResourceIPRecycleS) {
	malm.guard.Lock()
	defer malm.guard.Unlock()
	// 插入的队列
	heap.Push(&malm.timerHeap, &addrLeasePeriodTimerNodeS{
		record: record,
	})
	calm_utils.Debugf("push into timerHeap srv_instance_name[%s], k8sresource_id[%s] nspresource_release_time[%s] TotalCount[%d]",
		record.SrvInstanceName, record.K8SResourceID, record.NSPResourceReleaseTime.String(), len(malm.timerHeap))
	return
}

// 有两个入口，恢复、彻底删除
func (malm *mysqlAddrResourceLeasePeriodMgr) DelLeaseRecyclingRecord(k8sResourceID string) *table.TblK8SResourceIPRecycleS {
	malm.guard.Lock()
	defer malm.guard.Unlock()

	index := -1
	for _, timerNode := range malm.timerHeap {
		if timerNode.record.K8SResourceID == k8sResourceID {
			index = timerNode.index
			break
		}
	}

	if index != -1 {
		// 从最小堆中删除
		timerNode := malm.timerHeap[index]
		calm_utils.Debugf("remove from timerHeap srv_instance_name[%s], k8sresource_id[%s] nspresource_release_time[%s]",
			timerNode.record.SrvInstanceName, timerNode.record.K8SResourceID, timerNode.record.NSPResourceReleaseTime.String())

		heap.Remove(&malm.timerHeap, index)
		calm_utils.Debugf("timerHeap TotalCount[%d]", len(malm.timerHeap))
		return timerNode.record
	} else {
		calm_utils.Errorf("k8sResourceID[%s] not in timerHeap", k8sResourceID)
	}
	return nil
}

func (malm *mysqlAddrResourceLeasePeriodMgr) getMostRecentTimeout() (time.Time, int) {
	malm.guard.Lock()
	defer malm.guard.Unlock()

	if len(malm.timerHeap) > 0 {
		return malm.timerHeap[0].record.NSPResourceReleaseTime, len(malm.timerHeap)
	}
	// Time.IsZero() == true
	return time.Time{}, 0
}

func (malm *mysqlAddrResourceLeasePeriodMgr) timeCheckExpiration(now time.Time) {
	malm.guard.Lock()
	defer malm.guard.Unlock()

	// 这个已经是按时间排序的了
	var popCount int
	for _, timerNode := range malm.timerHeap {
		if now.After(timerNode.record.NSPResourceReleaseTime) {
			calm_utils.Debugf("Timeout srv_instance_name[%s], k8sresource_id[%s] nspresource_release_time[%s]",
				timerNode.record.SrvInstanceName, timerNode.record.K8SResourceID, timerNode.record.NSPResourceReleaseTime.String())
			// 调用db的删除数据
			malm.msm.expiredRecycling(timerNode.record)
			popCount++
		}
	}

	// 不要在slice range里面修改slice，range会赋值变量
	for popCount > 0 {
		heap.Pop(&malm.timerHeap)
		popCount--
	}
}

// NewAddrResourceLeasePeriodMgr 构造地址租期管理对象
func NewAddrResourceLeasePeriodMgr(ctx context.Context, msm *mysqlStoreMgr) AddrResourceLeasePeriodMgrItf {
	addrResourceLeasePeriodMgr := new(mysqlAddrResourceLeasePeriodMgr)
	addrResourceLeasePeriodMgr.ctx = ctx
	addrResourceLeasePeriodMgr.msm = msm
	addrResourceLeasePeriodMgr.timerHeap = make(addrLeasePeriodTimerHeap, 0)
	return addrResourceLeasePeriodMgr
}
