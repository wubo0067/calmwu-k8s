/*
 * @Author: calm.wu
 * @Date: 2020-07-16 15:28:56
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:49:49
 */

package remoteruntime

import "time"

// RuntimeService interface should be implemented by a container runtime.
type RuntimeService interface {
	// Close disconnect with runtime service
	Close()

	// ExecSync executes a command in the container, and returns the stdout output.
	// If command exits with a non-zero exit code, an error is returned.
	ExecSync(containerID string, cmd []string, timeout time.Duration) (data []byte, err error)

	RunBash(containerID string, cmd string) (data []byte, err error)
}
