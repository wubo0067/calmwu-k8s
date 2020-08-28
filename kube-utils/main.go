/*
 * @Author: calmwu
 * @Date: 2020-07-11 17:39:40
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-07-11 17:47:15
 */

package main

import (
	"log"
	"time"

	remoteruntime "kube-utils/remote_runtime"
)

func main() {
	runtimeService, err := remoteruntime.NewRemoteRuntimeService("", 3*time.Second)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer runtimeService.Close()

	data, err := runtimeService.RunBash("af1e3248721ea", "ls -al")
	if err != nil {
		log.Fatalf("RunBash failed, err:%s", err.Error())
		return
	}
	log.Printf("RunBash successed. response:%s\n", string(data))
}
