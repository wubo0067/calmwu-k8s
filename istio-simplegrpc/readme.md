### 1: 编译
*   make server  
*   make client  
*   make container  


### 2. 部署  
* Create configmap from descriptor set  
    `kubectl delete cm helloworld-proto-describe -n istio-ns`  
    `kubectl create configmap helloworld-proto-describe --from-file=helloworld.Greeter.pd -n istio-ns `

* Add annotations to your Deployment spec  
    `sidecar.istio.io/userVolume: '[{"name":"descriptor","configMap":{"name":"proto-descriptor","items":[{"key":"proto.pb","path":"proto.pb"}]}}]'`  
    `sidecar.istio.io/userVolumeMount: '[{"name":"descriptor","mountPath":"/etc/envoy"}]'`


### 3. 测试  

* get测试命令, 带上自定义头  
    ` curl -H "CallType: GRPC_Call" http://greeter.istio-ns.svc.cluster.local:8081/v1/say?name=sdsdsd -v `  


* post测试命令  
    ` curl -X POST -H "CallType: GRPC_Call" http://greeter.istio-ns.svc.cluster.local:8081/v1/reservations
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
