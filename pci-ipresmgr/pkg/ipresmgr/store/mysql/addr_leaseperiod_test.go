/*
 * @Author: calm.wu
 * @Date: 2019-09-02 17:35:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-02 19:45:13
 */

package mysql

import (
	"container/heap"
	"fmt"
	"pci-ipresmgr/table"
	"testing"
	"time"
)

func TestAddrTimerHeap(t *testing.T) {
	timerHeap := make(addrLeasePeriodTimerHeap, 0)
	heap.Init(&timerHeap)

	for i := 0; i < 10; i++ {
		heap.Push(&timerHeap, &addrLeasePeriodTimerNodeS{
			record: &table.TblK8SResourceIPRecycleS{
				K8SResourceID: fmt.Sprintf("inst-%d", i),
				NSPResourceReleaseTime: func() time.Time {
					if i%2 == 0 {
						return time.Now().Add(time.Duration(i) * time.Hour)
					} else {
						return time.Now().Add(time.Duration(-1*i) * time.Hour)
					}
				}(),
			},
		})
	}

	fmt.Printf("timeHeap len:%d\n", timerHeap.Len())

	for _, addrLeasePeriodTimerNode := range timerHeap {
		fmt.Printf("K8SResourceID:%s NSPResourceReleaseTime:%s index:%d\n",
			addrLeasePeriodTimerNode.record.K8SResourceID, addrLeasePeriodTimerNode.record.NSPResourceReleaseTime.String(),
			addrLeasePeriodTimerNode.index,
		)
	}

	fmt.Printf("------------------\n")

	// 删除一个指定的
	heap.Remove(&timerHeap, 3)
	for _, addrLeasePeriodTimerNode := range timerHeap {
		fmt.Printf("K8SResourceID:%s NSPResourceReleaseTime:%s index:%d\n",
			addrLeasePeriodTimerNode.record.K8SResourceID, addrLeasePeriodTimerNode.record.NSPResourceReleaseTime.String(),
			addrLeasePeriodTimerNode.index,
		)
	}

	fmt.Printf("------------------\n")

	for timerHeap.Len() > 0 {
		addrLeasePeriodTimerNode := heap.Pop(&timerHeap).(*addrLeasePeriodTimerNodeS)
		fmt.Printf("K8SResourceID:%s NSPResourceReleaseTime:%s index:%d\n",
			addrLeasePeriodTimerNode.record.K8SResourceID, addrLeasePeriodTimerNode.record.NSPResourceReleaseTime.String(),
			addrLeasePeriodTimerNode.index,
		)
	}

}
