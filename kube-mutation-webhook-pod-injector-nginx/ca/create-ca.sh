#!/bin/bash

set -v
set -x

rm *.pem *.csr 

#创建CA cert
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

#创建 webhook Server Cert -hostname={service-name}.{namespace}.svc
#
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -hostname=nginx-injector-pod-webhook-svc.nginx-injector-pod-webhook.svc \
    -profile=server webhook-server-csr.json | cfssljson -bare nginx-injector-pod-webhook-server