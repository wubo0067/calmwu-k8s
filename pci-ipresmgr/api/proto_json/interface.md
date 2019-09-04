### 1. 公共类型描述

#### 1.1 K8SApiResourceKindType
	const (
		K8SApiResourceKindDeployment K8SApiResourceKindType = iota // 0 = Deployment
		K8SApiResourceKindStatefulSet // 1 = StatefulSet
	)

### 2. WEBHOOK到IPResMgr接口描述

#### 2.1 创建IPPool

`Url: http://api.ipresmgr.com/v1/ippool/create`

`Method: Post`

`BodyType: Json`
```
type WB2IPResMgrCreateIPPoolReq struct {
	ReqID                  string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind     K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID           string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace           string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceID       string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"`
	K8SApiResourceName     string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	K8SApiResourceReplicas int                    `json:"K8SApiResourceReplicas" mapstructure:"K8SApiResourceReplicas"`
	NetRegionalID          string                 `json:"NetRegionalID" mapstructure:"NetRegionalID"`         // 网络域ID
	SubnetID               string                 `json:"SubnetID" mapstructure:"SubnetID"`                   // 子网ID
	SubnetGatewayAddr      string                 `json:"SubnetGatewayAddr" mapstructure:"SubnetGatewayAddr"` // 子网网关地址
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
	K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"`
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
}
```

#### 2.3 IPPool扩缩容

`Url: http://api.ipresmgr.com/v1/ippool/scale`

`Method: Post`

`BodyType: Json`
```
type WB2IPResMgrScaleIPPoolReq struct {
	ReqID                     string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind        K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID              string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace              string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceID          string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"`
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