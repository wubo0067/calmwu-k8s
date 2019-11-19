/*
 * @Author: calm.wu
 * @Date: 2019-10-04 11:38:28
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 12:02:17
 */

package protojson

// Maintain2IPResMgrRequestType 请求类型
type Maintain2IPResMgrRequestType int

const (
	//
	Maintain2IPResMgrDefaultReqType Maintain2IPResMgrRequestType = iota
	// Maintain2IPResMgrForceUnbindIP 解除pod和ip的绑定
	Maintain2IPResMgrForceUnbindIPReqType
	// Maintain2IPResMgrRequestReleaseIPPool 手工释放IPPool
	Maintain2IPResMgrForceReleaseK8SResourceIPPoolReqType
	// Maintain2IPResMgrForceReleaseK8SPodIP 手工释放Pod IP
	Maintain2IPResMgrForceReleaseK8SPodIPReqType
)

// Maintain2IPResMgrForceUnbindIPReq 运维系统强制解绑IP和POD请求命令
type Maintain2IPResMgrForceUnbindIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"` // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"` // 对应的namespace
	K8SPodName         string                 `json:"K8SPodName" mapstructure:"K8SPodName"`     // podname
}

// Maintain2IPResMgrForceReleaseK8SResourceIPReq 运维系统强制释放整个deployment、statfulset对应的IP资源，归还给nsp
type Maintain2IPResMgrForceReleaseK8SResourceIPPoolReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`             // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`             // 对应的namespace
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"` // deployment或statefulset的名字
}

// Maintain2IPResMgrForceReleasePodIPReq 运维系统强制释放IP，归还给nsp，针对单个pod
type Maintain2IPResMgrForceReleasePodIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`             // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`             // 对应的namespace
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"` // deployment或statefulset的名字
	K8SPodName         string                 `json:"K8SPodName" mapstructure:"K8SPodName"`                 // podname
}

// IPResMgr2MaintainRes 对Maintain的回应
type IPResMgr2MaintainRes struct {
	ReqID   string                       `json:"ReqID" mapstructure:"ReqID"`
	ReqType Maintain2IPResMgrRequestType `json:"ReqType" mapstructure:"ReqType"`
	Code    IPResMgrErrorCode            `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg     string                       `json:"Msg" mapstructure:"Msg"`
}
