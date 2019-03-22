创建deployment kubectl apply -f ng-as-deployment.yaml

进行扩容，kubectl scale deploy as-nginx-deployment --replicas 2

由于资源不够，有个pod状态时Pending
NAME                                   READY     STATUS    RESTARTS   AGE
as-nginx-deployment-2661823869-nmgvn   1/1       Running   0          1m
as-nginx-deployment-2661823869-pbz0n   0/1       Pending   0          14s

192.168.6.128:6443
curl https://10.254.0.1:6443 --cacert /run/secrets/kubernetes.io/serviceaccount/ca.crt

[root@localhost bin]# export KUBERNETES_SERVICE_HOST=192.168.6.128
[root@localhost bin]# export KUBERNETES_SERVICE_PORT=6443

token=`cat /run/secrets/kubernetes.io/serviceaccount/token`
curl https://10.254.0.1:443 --cacert /run/secrets/kubernetes.io/serviceaccount/ca.crt -H "Authorization: Bearer $token"

启动cluster_autoscaler
./cluster-autoscaler --kubeconfig=/home/calm/config --v=1

编译
go build -x -v -mod=vendor -o cluster-autoscaler main.go version.go