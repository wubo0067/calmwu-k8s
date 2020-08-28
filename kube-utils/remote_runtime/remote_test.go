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

	cmdLines := []string{
		"cat<<EOF\n",
		"123456\n",
		"654321\n",
		"EOF",
	}

	data, err := runtimeService.RunBash("af1e3248721ea", cmdLines)
	if err != nil {
		t.Errorf("RunBash failed, err:%s", err.Error())
		return
	}
	t.Logf("RunBash successed. response:%s\n", string(data))
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
