/*
 * @Author: calm.wu
 * @Date: 2019-07-10 15:00:31
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 14:51:15
 */

/*
Package protojson 协议
*/
package protojson

// K8SApiResourceKindType 资源类型
type K8SApiResourceKindType int16

const (
	// K8SApiResourceKindDeployment Deployment类型
	K8SApiResourceKindDeployment K8SApiResourceKindType = iota
	// K8SApiResourceKindStatefulSet Statefulset类型
	K8SApiResourceKindStatefulSet
	// K8SApiResourceKindJob Job类型
	K8SApiResourceKindJob
	// K8SApiResourceKindCronJob CronJob类型
	K8SApiResourceKindCronJob
	// K8SApiResourceKindUnknown
	K8SApiResourceKindUnknown
)

// WB2IPResMgrRequestType 请求类型
type WB2IPResMgrRequestType int

const (
	// WB2IPResMgrRequestCreateIPPool 创建IPPool
	WB2IPResMgrRequestCreateIPPool WB2IPResMgrRequestType = iota
	// WB2IPResMgrRequestReleaseIPPool 释放IPPool
	WB2IPResMgrRequestReleaseIPPool
	// WB2IPResMgrRequestScaleIPPool 扩缩容IPPool
	WB2IPResMgrRequestScaleIPPool
)

// WB2IPResMgrCreateIPPoolReq webhook层通知ipresmgr创建ippool
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
	SubnetCIDR             string                 `json:"SubnetCIDR" mapstructure:"SubnetCIDR"`                         // 子网CIDR，为了掩码
}

// WB2IPResMgrReleaseIPPoolReq webhook层通知ipresmgr释放ippool
type WB2IPResMgrReleaseIPPoolReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
}

// WB2IPResMgrScaleIPPoolReq webhook层通知ipresmgr修改副本数量，用户可以做update
type WB2IPResMgrScaleIPPoolReq struct {
	ReqID                     string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind        K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID              string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"`
	K8SNamespace              string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"`
	K8SApiResourceName        string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	K8SApiResourceOldReplicas int                    `json:"K8SApiResourceOldReplicas" mapstructure:"K8SApiResourceOldReplicas"`
	K8SApiResourceNewReplicas int                    `json:"K8SApiResourceNewReplicas" mapstructure:"K8SApiResourceNewReplicas"`
	NetRegionalID             string                 `json:"NetRegionalID" mapstructure:"NetRegionalID"`         // 网络域ID
	SubnetID                  string                 `json:"SubnetID" mapstructure:"SubnetID"`                   // 子网ID
	SubnetGatewayAddr         string                 `json:"SubnetGatewayAddr" mapstructure:"SubnetGatewayAddr"` // 子网网关地址
	SubnetCIDR                string                 `json:"SubnetCIDR" mapstructure:"SubnetCIDR"`               // 子网CIDR，为了掩码
}

// IPResMgr2WBRes 操作返回信息
type IPResMgr2WBRes struct {
	ReqID   string                 `json:"ReqID" mapstructure:"ReqID"`
	ReqType WB2IPResMgrRequestType `json:"ReqType" mapstructure:"ReqType"`
	Code    IPResMgrErrorCode      `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg     string                 `json:"Msg" mapstructure:"Msg"`
}
