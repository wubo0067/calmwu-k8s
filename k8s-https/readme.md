### 镜像打包
```
docker build -t littlebull/https-server:0.0.1 -f Dockerfile_server .
docker build -t littlebull/https-client:0.0.1 -f Dockerfile_client .

```

### 查看镜像结构
```
docker create --name https-srv littlebull/https-server:0.0.1
docker export -o https-srv.tar e6d3e3ebfbe6432097b87c9116fdcfd7c2929bd3267bcb41dbdf23fd0f6987ab
tar xf https-srv.tar
```


### 安装nginx ingress
```
https://www.qikqiak.com/post/install-nginx-ingress/
```

### curl访问
```
curl -v https://192.168.2.128/my-svc -H "host:k8s.calmwu.org" --cacert server.crt -k
```