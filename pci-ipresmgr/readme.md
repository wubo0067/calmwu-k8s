### 业务流程说明

#### 1. 预分配IP
`
WEBHOOK通知资源管理服务，deployment-name，副本数量，网络信息，服务收到请求后去nsp获取ip和网络相关资源。
获取资源后写入数据表中WB2IPResMgrCreateIPPoolReq，状态为未绑定。
`

#### 2. CNI获取IP
`
CNI通过接口，传入pod的标识，类型，服务从表中查询出可用的地址资源，返回给CNI
`

#### 3. CNI释放IP
`
CNI通过接口，传入pod标识，类型，服务将该ip状态设置为未绑定
curl -X POST http://192.168.6.134:30001/v1/ip/release
`

#### 4. 释放IP资源
`
WEBHOOK通知地址资源管理服务，删除对应的地址。服务响应请求，将回收的资源写入回收表中，设定租期，到期后将其返回给NSP
`

#### 5. 异常的情况，kubelet失效，节点被驱逐