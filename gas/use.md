### 启动consul
```
/usr/local/bin/consul agent -node node0 -bind 192.168.6.134 -dev -client 0.0.0.0
```

### 启动网关，web层使用的是api
```
pci-gateway --registry=consul api --handler=api --namespace=eci.v1.api

pci-gateway --registry=consul --transport=grpc api --handler=api --namespace=eci.v1.api
```

### 启动web层

web层是api实现的，example在micro-in-cn/all-in-one/basic-practices/micro-api
```
./web-hello --registry=consul

./web-hello --registry=consul --transport=grpc
```

### 观察网络，api和web层之间是长连接
```
tcp6       0      0 192.168.6.134:34391     192.168.6.134:33556     ESTABLISHED 22838/./api-hello   
tcp6       0      0 192.168.6.134:38791     192.168.6.134:54270     ESTABLISHED 23271/./api-hello
```

### 调用
- 使用get
```
curl "http://localhost:8080/hello/call?name=john"
[GET] Hello client john!
```
- 使用post
```
curl -H 'Content-Type: application/json' -d '{"name": "john"}' http://localhost:8080/hello/call
[POST] Hello client john!
```

### 启动srv层
所有模块都指定transport=grpc，这样才能打通三层
```
./srv-stringprocess --registry=consul --transport=grpc
./srv-usermgr --registry=consul --transport=grpc
```

### 配置管理使用consul集中配置
1. 在路径中设置配置信息

    consul kv put config/srv/k8soperator "$(cat ./config.json)"

### 启动调链跟踪Jaeger
```
docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.6
```
- 打开界面：http://192.168.6.134:16686/search