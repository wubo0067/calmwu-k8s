/*
 * @Author: calm.wu
 * @Date: 2020-07-16 15:09:55
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-07-16 20:50:38
 */

// Package remoteruntime is a gRPC implementation of internalapi.RuntimeService.
package remoteruntime

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/containerd/containerd/defaults"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	restclient "k8s.io/client-go/rest"
	remoteclient "k8s.io/client-go/tools/remotecommand"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/kubelet/util"
	utilexec "k8s.io/utils/exec"
)

type RemoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient runtimeapi.RuntimeServiceClient
	conn          *grpc.ClientConn
}

var _ RuntimeService = &RemoteRuntimeService{}

// NewRemoteRuntimeService creates a new RuntimeService.
func NewRemoteRuntimeService(endpoint string, connectionTimeout time.Duration) (RuntimeService, error) {
	endpoints := func() []string {
		if len(endpoint) > 0 {
			return []string{endpoint}
		}
		return _defaultRuntimeEndpoints
	}()

	conn, err := getConnection(endpoints, connectionTimeout)
	if err != nil {
		err = errors.Wrapf(err, "get connect failed.")
		klog.Error(err.Error())
		return nil, err
	}

	return &RemoteRuntimeService{
		timeout:       connectionTimeout,
		runtimeClient: runtimeapi.NewRuntimeServiceClient(conn),
		conn:          conn,
	}, nil
}

// getConnection connect to runtime server
func getConnection(endPoints []string, connectionTimeout time.Duration) (*grpc.ClientConn, error) {
	if len(endPoints) == 0 {
		return nil, fmt.Errorf("endpoint is not set")
	}

	endPointsLen := len(endPoints)
	var conn *grpc.ClientConn
	for indx, endPoint := range endPoints {
		klog.Infof("connect using endpoint '%s' with '%s' timeout", endPoint, connectionTimeout)
		addr, dialer, err := util.GetAddressAndDialer(endPoint)
		if err != nil {
			if indx == endPointsLen-1 {
				return nil, err
			}
			klog.Error(err)
			continue
		}
		conn, err = grpc.Dial(addr,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(connectionTimeout),
			grpc.WithContextDialer(dialer),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaults.DefaultMaxRecvMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaults.DefaultMaxSendMsgSize)),
		)
		if err != nil {
			errMsg := errors.Wrapf(err, "connect endpoint '%s', make sure you are running as root and the endpoint has been started", endPoint)
			if indx == endPointsLen-1 {
				return nil, errMsg
			}
			klog.Error(errMsg)
		} else {
			klog.V(3).Infof("connected successfully using endpoint: %s", endPoint)
			break
		}
	}
	return conn, nil
}

// Close disconnect with runtime service
func (r *RemoteRuntimeService) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

// ExecSync executes a command in the container, and returns the stdout output.
func (r *RemoteRuntimeService) ExecSync(containerID string, cmd []string, timeout time.Duration) (data []byte, err error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout != 0 {
		// Use timeout + default timeout (2 minutes) as timeout to leave some time for
		// the runtime to do cleanup.
		ctx, cancel = context.WithTimeout(context.Background(), r.timeout+timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	timeoutSeconds := int64(timeout.Seconds())
	req := &runtimeapi.ExecSyncRequest{
		ContainerId: containerID,
		Cmd:         cmd,
		Timeout:     timeoutSeconds,
	}
	resp, err := r.runtimeClient.ExecSync(ctx, req)
	if err != nil {
		klog.Errorf("ExecSync %s '%s' from runtime service failed: %v", containerID, strings.Join(cmd, " "), err)
		return nil, err
	}

	err = nil
	if resp.ExitCode != 0 {
		err = utilexec.CodeExitError{
			Err:  fmt.Errorf("command '%s' exited with %d: %s", strings.Join(cmd, " "), resp.ExitCode, resp.Stderr),
			Code: int(resp.ExitCode),
		}
	}

	return append(resp.Stdout, resp.Stderr...), err
}

func (r *RemoteRuntimeService) RunBash(containerID string, cmd string) ([]byte, error) {
	pr, pw := io.Pipe()

	request := &runtimeapi.ExecRequest{
		ContainerId: containerID,
		Cmd:         []string{"sh"},
		Tty:         false,
		Stdin:       true,
		Stdout:      true,
		Stderr:      true,
	}

	resp, err := r.runtimeClient.Exec(context.Background(), request)
	if err != nil {
		klog.Errorf("RunBash %s '/bin/sh' from runtime service failed: %v", containerID, err)
		return nil, err
	}

	execURL := resp.Url

	URL, err := url.Parse(execURL)
	if err != nil {
		klog.Errorf("RunBash %s url.Parse %s failed: %v", containerID, execURL, err)
		return nil, err
	}

	klog.Infof("RunBash URL: %v", URL)

	executor, err := remoteclient.NewSPDYExecutor(&restclient.Config{TLSClientConfig: restclient.TLSClientConfig{Insecure: true}}, "POST", URL)
	if err != nil {
		klog.Errorf("RunBash %s NewSPDYExecutor failed: %v", containerID, err)
		return nil, err
	}

	//_, stdout, stderr := dockerterm.StdStreams()
	writer := new(bytes.Buffer)
	streamOptions := remoteclient.StreamOptions{
		Stdout: writer,
		Stderr: writer,
		Tty:    false,
		Stdin:  pr,
	}

	klog.Infof("StreamOptions: %v", streamOptions)

	go func() {
		klog.Info("---stream start----")
		executor.Stream(streamOptions)
		klog.Info("---stream exit----")
	}()

	// 通过pipe写入命令
	len, err := pw.Write([]byte(cmd))
	if err != nil {
		klog.Errorf("RunBash %s pipe write failed: %v", containerID, err)
		return nil, err
	}

	// 回车执行命令
	pw.Write([]byte("\n"))
	klog.Infof("RunBash %s enter perform cmd, len:%d", containerID, len)

	// 获取结果
	time.Sleep(time.Second)
	cmdRes := new(bytes.Buffer)
	cmdRes.Write(writer.Bytes())
	klog.Infof("cmd res:%s", cmdRes.String())

	// 退出sh
	pw.Write([]byte("exit\n"))
	pw.Close()

	klog.Info("executor completed.")

	return cmdRes.Bytes(), nil
}
