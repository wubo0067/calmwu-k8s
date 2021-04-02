#!/bin/bash

#创建证书
umask 077; openssl genrsa -out calmwu.key 2048
openssl req -new -key calmwu.key -out calmwu.csr -subj "/CN=calmwu"
openssl x509 -req -in calmwu.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial -out calmwu.crt -days 3650 
openssl x509 -in calmwu.crt -text -noout

#将账户信息添加到k8s集群中
kubectl config set-credentials calmwu --client-certificate=./calmwu.crt --client-key=./calmwu.key --embed-certs=true

#创建账户，设置用户访问的集群
kubectl config set-context calmwu@kubernetes --cluster=kubernetes --user=calmwu

#查看账户
kubectl config get-contexts

#切换账户。当前账户没有权限访问api-server，需要serviceaccount、rolebinding，role
kubectl config use-context calmwu@kubernetes

#创建serviceaccount、rolebinding，role
kubectl apply -f serviceaccount.yaml
kubectl apply -f role.yaml
kubectl apply -f role-binding.yaml


