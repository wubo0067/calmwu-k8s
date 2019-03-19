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
systemctl restart kube-apiserver
systemctl restart kube-controller-manager
systemctl restart kube-scheduler
systemctl restart kubelet
systemctl restart kube-proxy

---------------------------------------------------------------------------

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

使用namespace查询
kubectl get pods --namespace=nm-nginxdeployment

#用yaml格式输出
kubectl get svc --namespace=nm-nginxdeployment -o yaml

kubectl get ep
--------------------------------------------------------------------------------
查询service的标签选择器
kubectl get svc mywebpodsvc -o jsonpath='{.spec.selector}'
可以根据标签进行查询
kubectl get pod calmwupod -o jsonpath='{.metadata.labels.app}'

kubectl get pods -l key1=value1,key2=value2

#一旦使用了namespace后必须带上该参数
kubectl get pods -l app=nginx --namespace=nm-nginxdeployment

查看标签
kubectl get pods --show-labels --namespace=nm-nginxdeployment
--------------------------------------------------------------------------------

登录到容器
docker exec -it 1470cfaa1b1c /bin/bash

root@myweb-d3560:/usr/local/tomcat# env |grep MYSQL_SERVICE
MYSQL_SERVICE_PORT=3306
MYSQL_SERVICE_HOST=mysql

最后数据库访问的问题
https://stackoverflow.com/questions/49204339/mysql-communications-link-failure-in-kubernetes-sample

Label Selector在kubernetes中重要的使用场景
1：RC上定义的Label Selector来筛选要监控的Pod副本的数量
2：kube-proxy进程通过Service的Label Selector来选择对应的Pod，自动建立起每个Service到对应Pod的请求转发路由表，从而实现Service的智能负债均衡机制
3：通过对某些Node定义特定的Label，并且在Pod定义文件中使用NodeSelector这种标签调度策略，实现Pod“定向调度”特性。

Replication controller
pod期待的副本数量(replicas)
用于筛选目标pod的label selector
当pod的副本数量小于预期数量的时候，用于创建新pod的pod模板

问题
1：什么是共享pod的ip，每个docker实例都有自己的ip地址，这个是挂在主机网桥上的

2：对于资源pod、ReplicationController区别在哪。资源之间的区别是什么？

--------------------------------------------------------------------------------
节点扩容
kubectl scale rc frontend --replicas=2
会在创建一个pod和svc对应上
Name:			frontend
Namespace:		default
Labels:			<none>
Selector:		tier=frontend
Type:			NodePort
IP:			10.254.110.244
Port:			<unset>	8889/TCP
NodePort:		<unset>	30003/TCP
Endpoints:		172.17.0.5:8080,172.17.0.6:8080
Session Affinity:	None
No events.
