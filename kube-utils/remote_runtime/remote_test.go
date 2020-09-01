/*
 * @Author: calm.wu
 * @Date: 2020-07-16 16:47:58
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:51:09
 */

package remoteruntime

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

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

	shell, err := runtimeService.NewBashShell("62ccc313a60fb")
	if err != nil {
		t.Errorf("NewBashShell failed, err:%s\n", err.Error())
		return
	}

	cmdLines := []string{
		"cat<<EOF\n",
		"123456\n",
		"654321\n",
		"calmwu\n",
		"EOF",
	}

	stdout, stderr, err := shell.ExecCmd(cmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			t.Errorf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	}

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
		"echo \"aaaaaaa\n",
		"bbbbbbwb\n",
		"\" |grep -n wb",
	}

	stdout, stderr, err = shell.ExecCmd(pipeCmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			t.Errorf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	}

	cmdLines = []string{
		"ls -al /dev",
	}

	stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			t.Errorf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	}

	cmdLines = []string{
		"uptime",
	}

	stdout, stderr, err = shell.ExecCmd(cmdLines, 1)
	if err != nil {
		if len(stderr) > 0 {
			t.Errorf("ExecCmd failed, stderr:\n%s\n", stderr)
		} else {
			t.Errorf("ExecCmd failed, %s\n", err.Error())
		}
	} else {
		t.Logf("ExecCmd successed. response:\n%s\n", stdout)
	}

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
