### 启动网关，web层使用的是api
```
micro --registry=consul api --handler=api --namespace=eci.v1.api

micro --registry=consul --transport=grpc api --handler=api --namespace=eci.v1.api
```

### 启动web层
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
#### 使用get
```
curl "http://localhost:8080/hello/call?name=john"
[GET] Hello client john!
```
#### 使用post
```
curl -H 'Content-Type: application/json' -d '{"name": "john"}' http://localhost:8080/hello/call
[POST] Hello client john!
```

### 启动srv层
所有模块都指定transport=grpc，这样才能打通三层
```
./srv-stringprocess --registry=consul --transport=grpc
```