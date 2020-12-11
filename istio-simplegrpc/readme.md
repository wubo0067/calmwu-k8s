### 1: 编译
*   make server  
*   make client  
*   make container  


### 2. 部署  
* Create configmap from descriptor set  
    `kubectl delete cm simplegrpcsrv-proto-describe -n istio-ns`  
    `kubectl create configmap simplegrpcsrv-proto-describe --from-file=istio-simplegrpc.pd -n istio-ns`

* Add annotations to your Deployment spec  
    `sidecar.istio.io/userVolume: '[{"name":"descriptor","configMap":{"name":"simplegrpcsrv-proto-describe","items":[{"key":"istio-simplegrpc.pd","path":"istio-simplegrpc.pd"}]}}]'`  
    `sidecar.istio.io/userVolumeMount: '[{"name":"descriptor","mountPath":"/etc/envoy"}]'`


### 3. 测试  

* get测试命令, 带上自定义头
    `curl -H "CallType: GRPC_Call" http://istio-simplegrpc.istio-ns.svc.cluster.local:8081/v1/say?name=sdsdsd -v`  


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

### 4. 绑定到网关  
* Create configmap from descriptor set  
    `kubectl delete cm simplegrpcsrv-proto-describe -n istio-system`  
    `kubectl create configmap simplegrpcsrv-proto-describe --from-file=istio-simplegrpc.pd -n istio-system`
