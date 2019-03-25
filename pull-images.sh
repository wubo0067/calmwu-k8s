#!/bin/bash
docker pull openstackmagnum/kubernetes-apiserver:v1.13.4
docker tag openstackmagnum/kubernetes-apiserver:v1.13.4 k8s.gcr.io/kube-apiserver:v1.13.4
docker rmi openstackmagnum/kubernetes-apiserver:v1.13.4

docker pull openstackmagnum/kubernetes-controller-manager:v1.13.4
docker tag openstackmagnum/kubernetes-controller-manager:v1.13.4 k8s.gcr.io/kube-controller-manager:v1.13.4
docker rmi openstackmagnum/kubernetes-controller-manager:v1.13.4

docker pull openstackmagnum/kubernetes-scheduler:v1.13.4
docker tag openstackmagnum/kubernetes-scheduler:v1.13.4 k8s.gcr.io/kube-scheduler:v1.13.4
docker rmi openstackmagnum/kubernetes-scheduler:v1.13.4

docker pull openstackmagnum/kubernetes-proxy:v1.13.4
docker tag openstackmagnum/kubernetes-proxy:v1.13.4 k8s.gcr.io/kube-proxy:1.2.6
docker rmi openstackmagnum/kubernetes-proxy:v1.13.4

docker pull openstackmagnum/pause:3.1
docker tag openstackmagnum/pause:3.1 k8s.gcr.io/pause:3.1
docker rmi openstackmagnum/pause:3.1

docker pull ibmcom/etcd:3.2.24
docker tag ibmcom/etcd:3.2.24 k8s.gcr.io/etcd:3.2.24
docker rmi ibmcom/etcd:3.2.24

docker pull coredns/coredns:1.2.6
docker tag coredns/coredns:1.2.6 k8s.gcr.io/coredns:1.2.6
docker rmi coredns/coredns:1.2.6

docker pull jmgao1983/flannel:v0.11.0-amd64
docker tag jmgao1983/flannel:v0.11.0-amd64 quay.io/coreos/flannel:v0.11.0-amd64
docker rmi jmgao1983/flannel:v0.11.0-amd64
