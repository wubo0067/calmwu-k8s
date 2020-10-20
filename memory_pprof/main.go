/*
 * @Author: CALM.WU
 * @Date: 2020-10-14 14:24:48
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-10-14 17:10:29
 */

package main

import (
	"fmt"
	"os"
	"time"

	calmUtils "github.com/wubo0067/calmwu-go/utils"
)

func allocate() []byte {
	return make([]byte, 1 << 20)
}


func main() {
	calmUtils.InstallPProf(30001)

	fmt.Println(os.Getpid())

	time.Sleep(5 * time.Second)

	var buf []byte

	fmt.Println("----------allocate start----------")

	for i := 1; i < 300; i++ {
		buf = append(buf, allocate()...)
		if i % 10 == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}

	fmt.Println("----------allocate end----------")
	select {}
}