/*
 * @Author: calm.wu
 * @Date: 2019-10-04 11:38:28
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 12:02:17
 */

package protojson

// Maintain2IPResMgrForceUnbindIPReq 运维系统强制解绑IP和POD请求命令
type Maintain2IPResMgrForceUnbindIPReq struct {
	ReqID              string                 `json:"ReqID" mapstructure:"ReqID"`
	K8SApiResourceKind K8SApiResourceKindType `json:"K8SApiResourceKind" mapstructure:"K8SApiResourceKind"`
	K8SClusterID       string                 `json:"K8SClusterID" mapstructure:"K8SClusterID"` // k8s集群id
	K8SNamespace       string                 `json:"K8SNamespace" mapstructure:"K8SNamespace"` // 对应的namespace
	K8SPodName         string                 `json:"K8SPodName" mapstructure:"K8SPodName"`     // podname
}

// Maintain2IPResMgrForceReleaseIPReq 运维系统强制释放IP，针对Job和CronJob
type Maintain2IPResMgrForceReleaseIPReq struct {
	ReqID        string `json:"ReqID" mapstructure:"ReqID"`
	K8SClusterID string `json:"K8SClusterID" mapstructure:"K8SClusterID"` // k8s集群id
	K8SNamespace string `json:"K8SNamespace" mapstructure:"K8SNamespace"` // 对应的namespace
	K8SPodName   string `json:"K8SPodName" mapstructure:"K8SPodName"`     // podname
}

// IPResMgr2MaintainRes 对Maintain的回应
type IPResMgr2MaintainRes struct {
	ReqID string            `json:"ReqID" mapstructure:"ReqID"`
	Code  IPResMgrErrorCode `json:"Code" mapstructure:"Code"` // 0 表示成功，!=0表示失败
	Msg   string            `json:"Msg" mapstructure:"Msg"`
}
