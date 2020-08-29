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

	cmdLines := []string{
		"ls -al",
	}

	shell, err := runtimeService.NewBashShell("62ccc313a60fb")
	if err != nil {
		log.Fatalf("NewBashShell failed, err:%s\n", err.Error())
		return
	}

	defer shell.Exit()

	stdout, stderr, err := shell.ExecCmd(cmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			log.Fatalf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			log.Fatalf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		log.Printf("ExecCmd successed. response:\n%s\n", stdout)
	}

	stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			log.Fatalf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			log.Fatalf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		log.Printf("ExecCmd successed. response:\n%s\n", stdout)
	}

	time.Sleep(100 * time.Second)
}
