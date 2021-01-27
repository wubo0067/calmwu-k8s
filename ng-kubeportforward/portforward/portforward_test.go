/*
 * @Author: CALM.WU
 * @Date: 2021-01-14 12:41:39
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-01-14 15:08:22
 */

package portforward

import (
	"sync"
	"testing"
	"time"
)

func TestGetFreePort(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			freePort, _ := GetFreePort()
			t.Logf("%d: port:%d", i, freePort)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestNilToChan(t *testing.T) {
	errCh := make(chan error)

	go func() {
		errCh <- func() error {
			return nil
		}()
		t.Log("insert nil to errCh")
	}()

	delayCh := time.After(3 * time.Second)

	select {
	case <-errCh:
		t.Log("receive nil error from errCh")
	case <-delayCh:
		t.Log("time out!")
	}
}
