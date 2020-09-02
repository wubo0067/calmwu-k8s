/*
 * @Author: calm.wu
 * @Date: 2020-09-02 11:13:01
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-02 12:45:03
 */

package workqueue

import (
	"log"
	"time"

	"k8s.io/client-go/util/workqueue"
)

func TestWorkQueueAdd() {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	item := "queue-item"
	queue.AddRateLimited(item)

	// 这个递增就是频率计数，不会对队列有任何行为
	// for i := 0; i < 1000; i++ {
	// 	queue.AddRateLimited(item)
	// }

	// 稍微等待
	time.Sleep(10 * time.Millisecond)

	log.Printf("AddRateLimited item:%s num requeues:%d", item, queue.NumRequeues(item))

	// Forget 将计数清零
	queue.Forget(item)
	log.Printf("Forget item:%s num requeues:%d", item, queue.NumRequeues(item))

	log.Printf("queue len:%d", queue.Len())

	// 获取，队列中只有一个
	for i := 0; i < queue.Len(); i++ {
		// 这个是从queue中获取数据，将对象插入processing set，从dirty set中删除
		gitem, _ := queue.Get()
		sitem := gitem.(string)
		log.Printf("get item:%s", sitem)
	}

	// 从processing set中删除，并且判断dirty set中是否存在，存在就插入队列末尾
	queue.Done(item)

	// 这里插入是有个时间差，如果duration=0，立即插入，否则会等待一段时间然后才插入队列
	queue.AddRateLimited("hello world")
	time.Sleep(time.Second)
	// 如果没有上面的等待，这里长度显示可能为0
	log.Printf("-----queue len:%d, num requeues:%d", queue.Len(), queue.NumRequeues(item))
	// Get是条件变量控制的，如果空就会一直等待的
	gitem, _ := queue.Get()
	sitem := gitem.(string)
	log.Printf("get item:%s queue len:%d, num requeues:%d", sitem, queue.Len(), queue.NumRequeues(item))
	queue.Forget(gitem)
	log.Printf("get item:%s queue len:%d, num requeues:%d", sitem, queue.Len(), queue.NumRequeues(item))

	queue.ShutDown()
}
