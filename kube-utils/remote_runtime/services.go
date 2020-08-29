/*
 * @Author: calm.wu
 * @Date: 2020-07-16 15:28:56
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:49:49
 */

package remoteruntime

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/klog"
)

type BashShell interface {
	// ExecCmd 执行命令，返回shell返回值和命令结果
	ExecCmd(cmdLines []string, timeOutSecs time.Duration) (string, string, error)

	// 结束shell
	Exit()
}

// RuntimeService interface should be implemented by a container runtime.
type RuntimeService interface {
	// Close disconnect with runtime service
	Close()

	// ExecSync executes a command in the container, and returns the stdout output.
	// If command exits with a non-zero exit code, an error is returned.
	ExecSync(containerID string, cmd []string, timeout time.Duration) (data []byte, err error)

	// NewBashShell 构造一个shell
	NewBashShell(containerID string) (BashShell, error)
}

type bashShellImpl struct {
	containerID string
	shWriter    *io.PipeWriter
	cmdStdout   *bytes.Buffer
	cmdStderr   *bytes.Buffer
	wg          sync.WaitGroup
}

var _ BashShell = &bashShellImpl{}

func (bsi *bashShellImpl) ExecCmd(cmdLines []string, timeOutSecs time.Duration) (string, string, error) {
	bsi.cmdStderr.Reset()
	bsi.cmdStdout.Reset()

	klog.Infof("cmdStdout len:%d", bsi.cmdStdout.Len())

	// 通过pipe写入命令
	writeLen := 0
	for _, cmdLine := range cmdLines {
		len, err := bsi.shWriter.Write(calmUtils.String2Bytes(cmdLine))
		if err != nil {
			err = errors.Wrapf(err, "ExecCmd %s write cmd failed.", bsi.containerID)
			klog.Errorf(err.Error())
			return "", "", err
		}
		writeLen += len
	}

	// 回车执行命令
	bsi.shWriter.Write(calmUtils.String2Bytes("\n"))
	klog.Infof("containerid[%s] ExecCmd write bytes:%d", bsi.containerID, writeLen)

	// 等待结果
	waitTimeout := time.Now().Add(timeOutSecs * time.Second)
	for {
		if bsi.cmdStdout.Len() > 0 || bsi.cmdStderr.Len() > 0 {
			klog.Infof("containerid[%s] ExecCmd cmdStdout len:%d, cmdStderr len:%d",
				bsi.containerID, bsi.cmdStdout.Len(), bsi.cmdStderr.Len())
			break
		}
		if time.Now().After(waitTimeout) {
			// 超时时间1秒
			err := errors.Errorf("ExecCmd %s timeout.", bsi.containerID)
			return "", "", err
		}

		time.Sleep(5 * time.Millisecond)
	}

	if bsi.cmdStderr.Len() > 0 {
		return "", bsi.cmdStderr.String(), errors.New("stderr")
	}

	return bsi.cmdStdout.String(), "", nil
}

func (bsi *bashShellImpl) Exit() {
	bsi.shWriter.Write(calmUtils.String2Bytes("exit\n"))
	bsi.shWriter.Close()
	bsi.wg.Wait()
}
