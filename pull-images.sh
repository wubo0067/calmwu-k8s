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


kubeadm join 192.168.6.131:6443 --token vawgip.ysfbx5kfogpb3u0o --discovery-token-ca-cert-hash sha256:43efc95b8d3edc43b3f3fdf6233395f593af48b8ac3a32c19bb7742bdf4b9160