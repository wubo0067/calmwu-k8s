/*
 * @Author: calm.wu
 * @Date: 2020-07-16 15:28:56
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:49:49
 */

package remoteruntime

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/klog"
)

type BashShell interface {
	// ExecCmd 执行命令，返回shell返回值和命令结果
	ExecCmd(cmdLines []string, timeoutDelay time.Duration) (string, error)

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
	containerID   string
	inPipeRead    *os.File
	inPipeWriter  *os.File
	outPipeRead   *os.File
	outPipeWriter *os.File
	wg            sync.WaitGroup
	guard         sync.Mutex
	// cmdStdout   *bytes.Buffer
	// cmdStderr   *bytes.Buffer
}

var _ BashShell = &bashShellImpl{}

// ExecCmd 执行命令，这个必须是串行，保证读取完整的命令返回
func (bsi *bashShellImpl) ExecCmd(cmdLines []string, timeoutDelay time.Duration) (string, error) {
	bsi.guard.Lock()
	defer bsi.guard.Unlock()

	// bsi.cmdStderr.Reset()
	// bsi.cmdStdout.Reset()

	//klog.Infof("cmdStdout len:%d", bsi.cmdStdout.Len())

	// 通过pipe写入命令
	writeLen := 0
	for _, cmdLine := range cmdLines {
		len, err := bsi.inPipeWriter.Write(calmUtils.String2Bytes(cmdLine))
		if err != nil {
			err = errors.Wrapf(err, "ExecCmd %s write cmd failed.", bsi.containerID)
			klog.Errorf(err.Error())
			return "", err
		}
		writeLen += len
	}

	// 回车执行命令
	bsi.inPipeWriter.Write(calmUtils.String2Bytes("\n"))
	klog.Infof("containerid[%s] ExecCmd write bytes:%d", bsi.containerID, writeLen)

	// 等待结果

	// 设置读取结果超时时间
	bsi.outPipeRead.SetDeadline(time.Now().Add(timeoutDelay))

	scanner := bufio.NewScanner(bsi.outPipeRead)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//fmt.Printf("data len:%d, atEof:%v\n", len(data), atEOF)
		if len(data) == 0 && atEOF {
			// 读到结束符
			return 0, nil, nil
		}

		if pos := strings.Index(calmUtils.Bytes2String(data), "0xEof"); pos >= 0 {
			//fmt.Printf("pos: %d\n", pos)
			// 返回偏移，自定义内容的长度，自定义数据的内容
			return pos + 5, data[0 : pos+5], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	for scanner.Scan() {
		// 读取到结果
		// fmt.Printf("read custom content:[%s]\n", scanner.Text())
		// TODO: 解析结果
		pos1 := bytes.LastIndexByte(scanner.Bytes(), '\n')
		pos2 := bytes.LastIndexByte(scanner.Bytes()[:pos1], '\n')
		resNum, err := strconv.Atoi(calmUtils.Bytes2String(scanner.Bytes()[pos2+1 : pos1]))

		if err != nil {
			return "", err
		}

		if resNum != 0 {
			err = errors.Errorf("cmd exec res code:%d", resNum)

			return "", err
		}

		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		err = errors.Wrap(err, "read cmd response failed.")

		return "", err
	}

	return "", nil
}

func (bsi *bashShellImpl) Exit() {
	bsi.inPipeWriter.Write(calmUtils.String2Bytes("exit\n"))
	bsi.inPipeWriter.Close()
	bsi.wg.Wait()
	bsi.inPipeRead.Close()
	bsi.outPipeWriter.Close()
	bsi.outPipeRead.Close()
}
