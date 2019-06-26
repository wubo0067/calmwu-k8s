### 启动网关，web层使用的是api
```
micro --registry=consul api --handler=api --namespace=eci.v1.api
```

### 启动web层
```
./api-hello --registry=consul
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

```
#### 使用post
```
curl -H 'Content-Type: application/json' -d '{"name": "john"}' http://localhost:8080/hello/call
```