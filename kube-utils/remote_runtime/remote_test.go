/*
 * @Author: calm.wu
 * @Date: 2020-07-16 16:47:58
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:51:09
 */

package remoteruntime

import (
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
