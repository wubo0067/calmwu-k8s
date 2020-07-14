/*
 * @Author: calmwu
 * @Date: 2020-07-11 17:40:27
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-07-11 17:45:27
 */

package main

import (
	"fmt"
	"net/url"
)

var (
	// DefaultRuntimeEndpoints unix系统默认的runtime地址
	DefaultRuntimeEndpoints = []string{"unix:///var/run/dockershim.sock", "unix:///run/containerd/containerd.sock", "unix:///run/crio/crio.sock"}
)

// ParseEndpoint 解析runtime endpoint
func ParseEndpoint(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}

	switch u.Scheme {
	case "tcp":
		return "tcp", u.Host, nil

	case "unix":
		return "unix", u.Path, nil

	case "":
		return "", "", fmt.Errorf("using %q as endpoint is deprecated, please consider using full url format", endpoint)

	default:
		return u.Scheme, "", fmt.Errorf("protocol %q not supported", u.Scheme)
	}
}
