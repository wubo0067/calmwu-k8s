/*
 * @Author: calm.wu
 * @Date: 2019-05-22 10:27:05
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-22 17:22:04
 */

package main

import (
	"crypto/tls"
	"flag"

	"k8s.io/klog"
)

type Config struct {
	CertFile string
	KeyFile  string
}

func (c *Config) addFlags() {
	flag.StringVar(&c.CertFile, "tls-cert-file", "/etc/kubernetes/pki/apiserver-etcd-client.crt", `
		File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated
		after server cert).`)
	flag.StringVar(&c.KeyFile, "tls-private-key-file", "/etc/kubernetes/pki/apiserver-etcd-client.key",
		"File containing the default x509 private key matching --tls-cert-file.")
}

func configTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}
