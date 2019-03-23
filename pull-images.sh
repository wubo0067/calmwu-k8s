#!/bin/bash
images=(kube-apiserver:v1.13.2 kube-controller-manager:v1.13.2 kube-scheduler:v1.13.2 kube-proxy:v1.13.2 pause:3.1 etcd:3.2.24)

for ima in ${images[@]}
do
   docker pull   mirrorgooglecontainers/$ima
   docker tag    mirrorgooglecontainers/$ima   k8s.gcr.io/$ima
   docker rmi  -f  mirrorgooglecontainers/$ima
done

docker pull coredns/coredns:1.2.6
docker tag coredns/coredns:1.2.6 k8s.gcr.io/coredns:1.2.6
docker rmi coredns/coredns:1.2.6