#!/bin/bash

kubectl delete secret nginx-injector-pod-webhook-certs -n nginx-injector-pod-webhook

#创建secret，作为证书的配置，部署文件中使用卷加载
kubectl create secret generic nginx-injector-pod-webhook-certs \
        --from-file=key.pem=../ca/nginx-injector-pod-webhook-server-key.pem \
        --from-file=cert.pem=../ca/nginx-injector-pod-webhook-server.pem \
        --dry-run=client -o yaml |
    kubectl -n nginx-injector-pod-webhook apply -f -