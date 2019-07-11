### 接口描述

#### 1. 创建IPPool

`Url: http://api.ipresmgr.com/v1/IPPool/Create`

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
	NetRegionalID          string                 `json:"NetRegionalID" mapstructure:"NetRegionalID"`
	SubnetID               string                 `json:"SubnetID" mapstructure:"SubnetID"`
}
```

#### 2. 释放IPPool

`Url: http://api.ipresmgr.com/v1/IPPool/Release`

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

### 3. IPPool扩缩容

`Url: http://api.ipresmgr.com/v1/IPPool/Scale`

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

### 4. 请求回应

`BodyType: Json`
```
type IPResMgr2WBRes struct {
	ReqID   string                 `json:"ReqID" mapstructure:"ReqID"`
	ReqType WB2IPResMgrRequestType `json:"ReqType" mapstructure:"ReqType"`
	Code    IPResMgrErrorCode      `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg     string                 `json:"Msg" mapstructure:"Msg"`
}
```