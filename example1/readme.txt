关闭防火墙

启动相关服务
systemctl enable etcd
systemctl enable docker
systemctl enable kube-apiserver
systemctl enable kube-controller-manager
systemctl enable kube-scheduler
systemctl enable kubelet
systemctl enable kube-proxy

systemctl start etcd
systemctl start docker
systemctl start kube-apiserver
systemctl start kube-controller-manager
systemctl start kube-scheduler
systemctl start kubelet
systemctl start kube-proxy

kubectl delete -f mysql-rc.yaml
kubectl delete -f mysql-svc.yaml

kubectl create -f mysql-rc.yaml
kubectl create -f mysql-svc.yaml

kubectl delete -f  myweb-svc.yaml
kubectl delete -f  myweb-rc.yaml

kubectl create -f  myweb-rc.yaml
kubectl create -f  myweb-svc.yaml

kubectl logs -f -c myweb myweb-d3560

kubectl describe svc myweb
kubectl describe pod myweb

kubectl get ep

查询service的标签选择器
kubectl get svc mywebpodsvc -o jsonpath='{.spec.selector}'

登录到容器
docker exec -it 1470cfaa1b1c /bin/bash

root@myweb-d3560:/usr/local/tomcat# env |grep MYSQL_SERVICE
MYSQL_SERVICE_PORT=3306
MYSQL_SERVICE_HOST=mysql

最后数据库访问的问题
https://stackoverflow.com/questions/49204339/mysql-communications-link-failure-in-kubernetes-sample



问题
1：什么是共享pod的ip，每个docker实例都有自己的ip地址，这个是挂在主机网桥上的

2：对于资源pod、ReplicationController区别在哪。资源之间的区别是什么？