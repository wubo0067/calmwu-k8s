/*
 * @Author: CALM.WU
 * @Date: 2021-01-14 11:35:21
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-01-14 15:45:04
 */

//
package portforward

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwardToPod struct {
	// restConfig is the kubernetes config
	restConfig *rest.Config
	// Pod is the selected pod for this port forwarding
	Pod v1.Pod
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int
	// PodPort is the target port for the pod
	PodPort int
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
	//
	pipeR, pipeW *os.File
}

func NewPortForward(config *rest.Config, podName string, nsName string, localPort int, podPort int, stopCh chan struct{}) *PortForwardToPod {
	r, w, _ := os.Pipe()

	return &PortForwardToPod{
		restConfig: config,
		Pod: v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podName,
				Namespace: nsName,
			},
		},
		LocalPort: localPort,
		PodPort:   podPort,
		StopCh:    stopCh,
		ReadyCh:   make(chan struct{}),
		Streams: genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    w,
			ErrOut: w,
		},
		pipeW: w,
		pipeR: r,
	}
}

func (pftp *PortForwardToPod) Start() error {

	errCh := make(chan error)

	go func() {
		path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
			pftp.Pod.ObjectMeta.Namespace, pftp.Pod.ObjectMeta.Name)
		hostIP := strings.TrimLeft(pftp.restConfig.Host, "htps:/")

		transport, upgrader, err := spdy.RoundTripperFor(pftp.restConfig)
		if err != nil {
			err = errors.Wrap(err, "spdy RoundTripperFor failed")
			calmUtils.Error(err.Error())
			errCh <- err
			return
		}

		dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
		fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", pftp.LocalPort, pftp.PodPort)}, pftp.StopCh, pftp.ReadyCh, pftp.Streams.Out, pftp.Streams.ErrOut)
		if err != nil {
			err = errors.Wrap(err, "portforward New failed")
			calmUtils.Error(err.Error())
			errCh <- err
			return
		}

		err = fw.ForwardPorts()
		if err != nil {
			err = errors.Wrapf(err, "%s forward ports failed", path)
			calmUtils.Error(err.Error())
			errCh <- err
		} else {
			calmUtils.Debugf("%s exit", path)
		}
	}()

	// 读取信息
	go func() {
		scanner := bufio.NewScanner(pftp.pipeR)
		for scanner.Scan() {
			calmUtils.Debugf("---%s---", scanner.Text())
		}
		calmUtils.Debug("log portForward exit!")
	}()

	select {
	case err := <-errCh:
		return err
	case <-pftp.ReadyCh:
		calmUtils.Debugf("Port forwarding is ready to get traffic. have fun!")
	}

	return nil
}

func (pftp *PortForwardToPod) Stop() {
	if pftp != nil {
		close(pftp.StopCh)
		pftp.pipeR.Close()
		pftp.pipeW.Close()
	}
}

// GetFreePort 获得一个本地空闲端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
