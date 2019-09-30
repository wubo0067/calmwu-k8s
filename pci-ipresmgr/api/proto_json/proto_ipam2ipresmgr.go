/*
 * @Author: calm.wu
 * @Date: 2019-07-11 14:37:19
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-11 16:01:01
 */

package protojson

// IPAM2IPResMgrRequireIPReq ipam向ipresmgr请求ip地址
type IPAM2IPResMgrRequireIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"` // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"` // 对应的namespace
	K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"` // deployment、statefulset、job、cronjob 名字
	K8SPodName         string                 `json:"K8SPodName" mapstructure:"K8SPodName"` // podname
	//K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"` // Deployment-id 或 StatefulSet-id
	//K8SPodID           string                 `json:"K8SPodID" mapstructure:"K8SPodID"` // pod-id 不是podname
}

// IPResMgr2IPAMRequireIPRes ipresmgr给ipam返回ip地址
type IPResMgr2IPAMRequireIPRes struct {
	ReqID             string            `json:"ReqID" mapstructure:"ReqID"`
	IP                string            `json:"IP" mapstructure:"IP"`
	MacAddr           string            `json:"MacAddr" mapstructure:"MacAddr"`
	PortID            string            `json:"PortID" mapstructure:"PortID"`
	SubnetGatewayAddr string            `json:"SubnetGatewayAddr" mapstructure:"SubnetGatewayAddr"`
	Code              IPResMgrErrorCode `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg               string            `json:"Msg" mapstructure:"Msg"`
}

// IPAM2IPResMgrReleaseIPReq ipam向ipresmgr释放ip地址
type IPAM2IPResMgrReleaseIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"` 
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"` // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"` // 对应的namespace
	K8SPodName         string                 `json:"K8SPodName" mapstructure:"K8SPodName"`     // podname
	IP                 string                 `json:"IP" mapstructure:"IP"`
	//K8SApiResourceID   string                 `json:"K8SApiResourceID" mapstructure:"K8SApiResourceID"` // Deployment-id 或 StatefulSet-id
	//K8SApiResourceName string                 `json:"K8SApiResourceName" mapstructure:"K8SApiResourceName"`
	//K8SPodID           string                 `json:"K8SPodID" mapstructure:"K8SPodID"` // pod-id 不是podname
}

// IPResMgr2IPAMReleaseIPRes 回应包
type IPResMgr2IPAMReleaseIPRes struct {
	ReqID string            `json:"ReqID" mapstructure:"ReqID"`
	Code  IPResMgrErrorCode `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg   string            `json:"Msg" mapstructure:"Msg"`
}
