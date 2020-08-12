/*
 * @Author: calmwu
 * @Date: 2020-08-01 13:07:10
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-08-01 13:55:12
 */

package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile("../../cert/server.crt")
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile("../../cert/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "0.0.0.0",
	}
}

func main() {
	server := &http.Server{
		Addr:         "0.0.0.0:8443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig(),
	}

	//// Having this does not change anything but just showing.
	//// go get -u golang.org/x/net/http2
	//if err := http2.ConfigureServer(server, nil); err != nil {
	//	log.Fatal(err)
	//}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("recv https request from:%s\n", r.RemoteAddr)
		w.Write([]byte(fmt.Sprintf("Protocol: %s", r.Proto)))
	})

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
