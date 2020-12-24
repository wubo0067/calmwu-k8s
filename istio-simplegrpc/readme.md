### 1: 编译
-------

*   make server  
*   make client  
*   make container  


### 2. 部署 
------- 
* Create configmap from descriptor set  
    `kubectl delete cm simplegrpcsrv-proto-describe -n istio-ns`  
    `kubectl create configmap simplegrpcsrv-proto-describe --from-file=istio-simplegrpc.pd -n istio-ns`

* Add annotations to your Deployment spec  
    `sidecar.istio.io/userVolume: '[{"name":"descriptor","configMap":{"name":"simplegrpcsrv-proto-describe","items":[{"key":"istio-simplegrpc.pd","path":"istio-simplegrpc.pd"}]}}]'`  
    `sidecar.istio.io/userVolumeMount: '[{"name":"descriptor","mountPath":"/etc/envoy"}]'`

* 启动服务  
   `kubectl apply -f deploy_server.yaml`  

* 修改istio-proxy loglevel  
    `kubectl exec -it istio-simplegrpc-server-v1-75785b9d7d-dm47b -c istio-proxy -n istio-ns -- curl -X POST localhost:15000/logging?level=debug`


### 3. 测试  
-------
* get测试命令, 带上自定义头，在pod内部访问  
    `curl -H "CallType: GRPC_Call" http://istio-simplegrpc.istio-ns.svc.cluster.local:8081/v1/say?name=sdsdsd -v`   
    `curl -H "CallType: GRPC_Echo" http://istio-simplegrpc.istio-ns.svc.cluster.local:8081/v1/echotimeout?message=sdsdsd -v`  

* 通过nodeport访问  
    `curl -H 'Host:www.istio-simplegrpc.com' -H 'CallType:GRPC_Call' http://192.168.6.128:32197/v1/say?name=sdsdsd -v`   


* post测试命令  
    ` curl -X POST -H "CallType: GRPC_Call" http://istio-simplegrpc.istio-ns.svc.cluster.local:8081/v1/reservations
    -d '{
        "title": "Lunchmeeting",
        "venue": "JDriven Coltbaan 3",
        "room": "atrium",
        "timestamp": "2018-10-10T11:12:13",
        "attendees": [
        {
        "ssn": "1234567890",
        "firstName": "Jimmy",
        "lastName": "Jones"
        },
        {
        "ssn": "9999999999",
        "firstName": "Dennis",
        "lastName": "Richie"
        }
        ]
        }' `

### 4. 使用configmap，绑定到网关  
-------
* Create configmap from descriptor set  
    `kubectl delete cm simplegrpcsrv-proto-describe -n istio-system`  
    `kubectl create configmap simplegrpcsrv-proto-describe --from-file=istio-simplegrpc.pd -n istio-system`

### 5. 从网关外部访问  
-------
* 编写网关和virtualservice资源，simplegrpc-gateway.yaml  
* 访问命令, 这里指定了Host。`curl -H 'Host:www.istio-simplegrpc.com' -H 'CallType:GRPC_Call' http://192.168.6.128:32197/v1/say?name=sdsdsd -v` 
* 结果：  
```
    [root@localhost networking]# curl -H 'Host:www.istio-simplegrpc.com' -H 'CallType:GRPC_Call1' http://192.168.6.128:32197/v1/say?name=sdsdsd -v
    *   Trying 192.168.6.128...
    * TCP_NODELAY set
    * Connected to 192.168.6.128 (192.168.6.128) port 32197 (#0)
    > GET /v1/say?name=sdsdsd HTTP/1.1
    > Host:www.istio-simplegrpc.com
    > User-Agent: curl/7.61.1
    > Accept: */*
    > CallType:GRPC_Call1
    > 
    < HTTP/1.1 200 OK
    < content-type: application/json
    < x-envoy-upstream-service-time: 4
    < grpc-status: 0
    < grpc-message: 
    < content-length: 92
    < date: Mon, 14 Dec 2020 08:23:15 GMT
    < server: istio-envoy
    < 
    {
    "message": "srv-host:istio-simplegrpc-server-v1-75785b9d7d-j88rt index:7 Hello sdsdsd"
    }
```   
