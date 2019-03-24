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


kubeadm join 192.168.6.131:6443 --token vawgip.ysfbx5kfogpb3u0o --discovery-token-ca-cert-hash sha256:43efc95b8d3edc43b3f3fdf6233395f593af48b8ac3a32c19bb7742bdf4b9160 --ignore-preflight-errors=Swap

vim /etc/docker/daemon.conf 加入{"registry-mirrors": ["http://f1361db2.m.daocloud.io"]}

curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://f1361db2.m.daocloud.io

我把kube-flannel.yml依赖的镜像版本从v0.11改为v0.10，我从这里下载的，docker pull cnych/flannel:v0.10.0-amd64


[root@k8snode2 calmwu]# docker pull mirrorgooglecontainers/pause:3.1
3.1: Pulling from mirrorgooglecontainers/pause
67ddbfb20a22: Pull complete 
Digest: sha256:59eec8837a4d942cc19a52b8c09ea75121acc38114a2c68b98983ce9356b8610
Status: Downloaded newer image for mirrorgooglecontainers/pause:3.1
[root@k8snode2 calmwu]# docker tag mirrorgooglecontainers/pause:3.1 k8s.gcr.io/pause:3.1
[root@k8snode2 calmwu]# docker rmi mirrorgooglecontainers/pause:3.1

docker pull mirrorgooglecontainers/kube-proxy:v1.13.2
docker tag mirrorgooglecontainers/kube-proxy:v1.13.2 k8s.gcr.io/kube-proxy:v1.13.2
docker rmi mirrorgooglecontainers/kube-proxy:v1.13.2

如果子节点起不来就重启下该服务
systemctl restart kubelet
