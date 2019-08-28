/*
 * @Author: calm.wu
 * @Date: 2019-07-10 15:00:31
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-11 14:51:39
 */

package protojson

type K8SApiResourceKindType int16

const (
	K8SApiResourceKindDeployment K8SApiResourceKindType = iota
	K8SApiResourceKindStatefulSet
)

type WB2IPResMgrRequestType int

const (
	WB2IPResMgrRequestCreateIPPool WB2IPResMgrRequestType = iota
	WB2IPResMgrRequestReleaseIPPool
	WB2IPResMgrRequestScaleIPPool
)

// webhook层通知ipresmgr创建ippool
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

// webhook层通知ipresmgr释放ippool
type WB2IPResMgrReleaseIPPoolReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"`
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
}

// webhook层通知ipresmgr扩缩容ippool
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

// 操作返回信息
type IPResMgr2WBRes struct {
	ReqID   string                 `json:"ReqID" mapstructure:"ReqID"`
	ReqType WB2IPResMgrRequestType `json:"ReqType" mapstructure:"ReqType"`
	Code    IPResMgrErrorCode      `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg     string                 `json:"Msg" mapstructure:"Msg"`
}
