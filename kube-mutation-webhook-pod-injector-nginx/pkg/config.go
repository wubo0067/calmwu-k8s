/*
 * @Author: CALM.WU
 * @Date: 2021-04-29 14:26:23
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-29 14:28:40
 */

// Package pkg is implement nginx injector to pod
package pkg

// SvrParamenters is server parameters
type SvrParamenters struct {
	Port           int    // webhook server port
	CertFile       string // path to the x509 certificate for https
	KeyFile        string // path to the x509 private key matching `CertFile`
	SidecarCfgFile string // path to sidecar injector configuration file
}
