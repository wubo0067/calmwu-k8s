### 调用
```
curl -H 'Content-Type: application/json' -d '{"name": "john"}' "http://localhost:8080/NamespaceSvr/GetNamespace"
```

### 启动rpc网关，自定义的namespace
```
micro --registry=consul --transport=grpc api --namespace=eci.v1.api --handler=rpc
```

### 启动服务
```
./main_eci_namespace --registry=consul --transport=grpc
```