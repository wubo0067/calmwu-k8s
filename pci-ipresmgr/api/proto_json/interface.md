### 1. 公共类型描述

#### 1.1 K8SApiResourceKindType
const (
	// K8SApiResourceKindDeployment Deployment类型
	K8SApiResourceKindDeployment K8SApiResourceKindType = iota
	// K8SApiResourceKindStatefulSet Statefulset类型
	K8SApiResourceKindStatefulSet
	// K8SApiResourceKindJob Job类型
	K8SApiResourceKindJob
	// K8SApiResourceKindCronJob CronJob类型
	K8SApiResourceKindCronJob
)

### 2. WEBHOOK到IPResMgr接口描述

#### 2.1 创建IPPool

`Url: http://api.ipresmgr.com/v1/ippool/create`

`Method: Post`

`BodyType: Json`
```
type WB2IPResMgrCreateIPPoolReq struct {
	ReqID                  string                 `json:"ReqID" mapstructure:"ReqID"`                                   // 请求id，消息对应
	K8SApiResourceKind     K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`         // workload的类型
	K8SClusterID           string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`                     // 集群标识
	K8SNamespace           string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`                     // 名字空间
	K8SApiResourceName     string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`         // deployment或statefulset的名字
	K8SApiResourceReplicas int                    `json:"K8SApiResourceReplicas" mapstructure:"K8SApiResourceReplicas"` // 副本数量
	NetRegionalID          string                 `json:"NetRegionalID" mapstructure:"NetRegionalID"`                   // 网络域ID
	SubnetID               string                 `json:"SubnetID" mapstructure:"SubnetID"`                             // 子网ID
	SubnetGatewayAddr      string                 `json:"SubnetGatewayAddr" mapstructure:"SubnetGatewayAddr"`           // 子网网关地址
}
```

#### 2.2 释放IPPool

`Url: http://api.ipresmgr.com/v1/ippool/release`

`Method: Post`

`BodyType: Json`
```
type WB2IPResMgrReleaseIPPoolReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
}
```

#### 2.3 IPPool扩缩容，webhook层通知ipresmgr修改副本数量，用户可以做update

`Url: http://api.ipresmgr.com/v1/ippool/scale`

`Method: Post`

`BodyType: Json`
```
type WB2IPResMgrScaleIPPoolReq struct {
	ReqID                     string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind        K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID              string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace              string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceName        string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	K8SApiResourceOldReplicas int                    `json:"K8SApiResourceOldReplicas" mapstructure:"K8SApiResourceOldReplicas"`
	K8SApiResourceNewReplicas int                    `json:"K8SApiResourceNewReplicas" mapstructure:"K8SApiResourceNewReplicas"`
}
```

#### 2.4 请求回应

`BodyType: Json`
```
type IPResMgr2WBRes struct {
	ReqID   string                 `json:"ReqID" mapstructure:"ReqID"`
	ReqType WB2IPResMgrRequestType `json:"ReqType" mapstructure:"ReqType"`
	Code    IPResMgrErrorCode      `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg     string                 `json:"Msg" mapstructure:"Msg"`
}
```

### 3. IPAM到IPResMgr接口描述

#### 3.1 IPAM获取IP

`Url: http://api.ipresmgr.com/v1/ip/require`

`Method: Post`

`BodyType: Json`
```
// ipam向ipresmgr请求ip地址
type IPAM2IPResMgrRequireIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`         // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`         // 对应的namespace
	K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"` // Deployment-id 或 StatefulSet-id
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	K8SPodID           string                 `json:"K8SPodID" mapstructure:"K8SPodID"` // pod-id 不是podname
}

// ipresmgr给ipam返回ip地址
type IPResMgr2IPAMRequireIPRes struct {
	ReqID             string            `json:"ReqID" mapstructure:"ReqID"`
	IP                string            `json:"IP" mapstructure:"IP"`
	MacAddr           string            `json:"MacAddr" mapstructure:"MacAddr"`
	PortID            string            `json:"PortID" mapstructure:"PortID"`
	SubnetGatewayAddr string            `json:"SubnetGatewayAddr" mapstructure:"SubnetGatewayAddr"`
	Code              IPResMgrErrorCode `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg               string            `json:"Msg" mapstructure:"Msg"`
}
```

#### 3.2 IPAM释放IP
`Url: http://api.ipresmgr.com/v1/ip/release`

`Method: Post`

`BodyType: Json`
```
type IPAM2IPResMgrReleaseIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`         // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`         // 对应的namespace
	K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"` // Deployment-id 或 StatefulSet-id
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	K8SPodID           string                 `json:"K8SPodID" mapstructure:"K8SPodID"`  // pod-id 不是podname
	IP                 string                 `json:"IP" mapstructure:"IP"`
}

type IPResMgr2IPAMReleaseIPRes struct {
	ReqID string            `json:"ReqID" mapstructure:"ReqID"`
	Code  IPResMgrErrorCode `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg   string            `json:"Msg" mapstructure:"Msg"`
}
```