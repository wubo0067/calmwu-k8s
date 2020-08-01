/*
 * @Author: calmwu
 * @Date: 2020-08-01 13:44:56
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-08-01 13:55:13
 */

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	client := &http.Client{Transport: transport2()}

	res, err := client.Get("https://my-svc.calm-space.svc:8443")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	fmt.Printf("Code: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", body)
}

func transport2() *http2.Transport {
	return &http2.Transport{
		TLSClientConfig:    tlsConfig(),
		DisableCompression: true,
		AllowHTTP:          false,
	}
}

//func transport1() *http.Transport {
//	return &http.Transport{
//		// Original configurations from `http.DefaultTransport` variable.
//		DialContext: (&net.Dialer{
//			Timeout:   30 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).DialContext,
//		ForceAttemptHTTP2:     true, // Set it to false to enforce HTTP/1
//		MaxIdleConns:          100,
//		IdleConnTimeout:       90 * time.Second,
//		TLSHandshakeTimeout:   10 * time.Second,
//		ExpectContinueTimeout: 1 * time.Second,
//
//		// Our custom configurations.
//		ResponseHeaderTimeout: 10 * time.Second,
//		DisableCompression:    true,
//		// Set DisableKeepAlives to true when using HTTP/1 otherwise it will cause error: dial tcp [::1]:8090: socket: too many open files
//		DisableKeepAlives:     false,
//		TLSClientConfig:       tlsConfig(),
//	}
//}

func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile("../../cert/server.crt")
	if err != nil {
		log.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false, /*InsecureSkipVerify用来控制客户端是否证书和服务器主机名。如果设置为true,则不会校验证书以及证书中的主机名和服务器主机名是否一致。
		因为在我们的例子中使用自签名的证书，所以设置它为true,仅仅用于测试目的。*/
		//ServerName:         "localhost",
		ServerName: "my-svc.calm-space.svc.cluster.local",
	}
}
