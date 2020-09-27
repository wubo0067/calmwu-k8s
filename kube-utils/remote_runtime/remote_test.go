/*
 * @Author: calm.wu
 * @Date: 2020-07-16 16:47:58
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:51:09
 */

package remoteruntime

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"k8s.io/utils/exec"
)

func TestExecSync(t *testing.T) {
	runtimeService, err := NewRemoteRuntimeService("", 3*time.Second)
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer runtimeService.Close()

	data, err := runtimeService.ExecSync("72d13d5d3c475", []string{"ls", "-al"}, 3*time.Second)
	if err != nil {
		exit, ok := err.(exec.ExitError)
		if ok {
			if exit.ExitStatus() == 0 {
				t.Logf("exec successed, procsss exit code is 0. response:%s\n", string(data))
			}
		}
		t.Errorf("exec failed, err:%s", err.Error())
		return
	}
	t.Logf("exec successed. response:%s\n", string(data))
}

func TestRunBash(t *testing.T) {
	runtimeService, err := NewRemoteRuntimeService("", 3*time.Second)
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer runtimeService.Close()

	shell, err := runtimeService.NewBashShell("ceff738df4c57")
	if err != nil {
		t.Errorf("NewBashShell failed, err:%s\n", err.Error())
		return
	}

	// cmdLines := []string{
	// 	"ls -al /dev;echo $?;echo 0xEof",
	// }

	// out, err := shell.ExecCmd(cmdLines, 200*time.Millisecond)
	// if err != nil {
	// 	t.Errorf("ExecCmd failed, %s\n", err.Error())
	// } else {
	// 	t.Logf("ExecCmd successed. response:\n%s\n", out)
	// }

	// time.Sleep(1 * time.Second)

	// stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	// if err != nil {
	// 	if len(stderr) > 0 {
	// 		t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
	// 	} else {
	// 		t.Errorf("ExecCmd failed, %s\n", err.Error())
	// 	}
	// } else {
	// 	t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	// }

	// time.Sleep(1 * time.Second)

	// stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	// if err != nil {
	// 	if len(stderr) > 0 {
	// 		t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
	// 	} else {
	// 		t.Errorf("ExecCmd failed, %s\n", err.Error())
	// 	}
	// } else {
	// 	t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	// }

	pipeCmdLines := []string{
		"echo bbbbbbbbwb|grep -n wb;echo $?;echo 0xEof",
	}

	_, err = shell.ExecCmd(pipeCmdLines, 200*time.Millisecond)
	if err != nil {
		t.Errorf("ExecCmd failed, %s\n", err.Error())
	}

	pipeCmdLines = []string{
		"sdsd bbbbbbbbwb|grep -n wb;echo $?;echo 0xEof",
	}

	_, err = shell.ExecCmd(pipeCmdLines, 200*time.Millisecond)
	if err != nil {
		t.Errorf("ExecCmd failed, %s\n", err.Error())
	}

	// cmdLines = []string{
	// 	"ls -al /dev",
	// }

	// stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	// if err != nil {
	// 	if len(stderr) > 0 {
	// 		t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
	// 	} else {
	// 		t.Errorf("ExecCmd failed, %s\n", err.Error())
	// 	}
	// } else {
	// 	t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	// }

	// cmdLines = []string{
	// 	"uptime",
	// }

	// stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	// if err != nil {
	// 	if len(stderr) > 0 {
	// 		t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
	// 	} else {
	// 		t.Errorf("ExecCmd failed, %s\n", err.Error())
	// 	}
	// } else {
	// 	t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	// }

	shell.Exit()
}

func TestMulitiExec(t *testing.T) {
	runtimeService, err := NewRemoteRuntimeService("", 3*time.Second)
	if err != nil {
		t.Error(err.Error())

		return
	}

	defer runtimeService.Close()

	pipeCmdLines := []string{
		"echo \"aaaaaaa\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"wbbbbbdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddbwb\n",
		"\" |grep -n wb;echo $?;echo 0xEof",
	}

	shell, err := runtimeService.NewBashShell("ceff738df4c57")
	if err != nil {
		t.Errorf("NewBashShell failed, err:%s\n", err.Error())

		return
	}

	timeStart := time.Now()

	for i := 0; i < 200; i++ {
		shell.ExecCmd(pipeCmdLines, 200*time.Millisecond)
	}

	t.Logf("execute consume time: %s", time.Since(timeStart).String())

	shell.Exit()
}

type readerWrapper struct {
	reader io.Reader
}

func (r readerWrapper) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func TestOSPipe(t *testing.T) {
	pr, pw := io.Pipe()

	go func() {
		io.Copy(os.Stdout, readerWrapper{pr})
		fmt.Println("Pipe read EOF")
	}()

	pw.Write([]byte("---Hello world---\n"))
	time.Sleep(time.Second)

	fmt.Println("Pipe close")
	pw.Close()

	time.Sleep(time.Second)
	fmt.Println("TestOSPipe exit")
}

func TestPipeScan(t *testing.T) {
	pr, pw, _ := os.Pipe()

	go func() {
		pr.SetDeadline(time.Now().Add(time.Second))

		scanner := bufio.NewScanner(pr)
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
			fmt.Printf("read custom content:[%s]\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("scanner error! read pipe failed. err:%s\n", err.Error())
			return
		}

		fmt.Printf("scanner pipe routine exit\n")
	}()

	_, _ = pw.Write([]byte("---Hello world---\n"))
	_, _ = pw.Write([]byte("0xEof"))
	_, _ = pw.Write([]byte("---Hello sci---\n"))
	_, _ = pw.Write([]byte("0xEof"))

	toC := time.After(3 * time.Second)
	// 读取超时，就关闭pipe
	<-toC
	fmt.Println("close pipe")
	pw.Close()

	time.Sleep(2 * time.Second)
	pr.Close()
	fmt.Println("TestScanFromPipe exit")
}
