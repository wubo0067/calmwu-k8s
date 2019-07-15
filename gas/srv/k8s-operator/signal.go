// +build linux
/*
 * @Author: calm.wu
 * @Date: 2019-07-15 16:43:48
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-15 16:51:11
 */

package k8soperator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigCh
			switch sig {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				cancel()
				return
			case syscall.SIGUSR1:
				fallthrough
			case syscall.SIGUSR2:
				calm_utils.DumpStacks()
			}
		}
	}()
}
