/*
 * @Author: calm.wu
 * @Date: 2019-08-28 15:25:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 23:49:07
 */

package srv

import (
	"net/http"

	proto "pci-ipresmgr/api/proto_json"

	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func cniRequireIP(c *gin.Context) {
	var req proto.IPAM2IPResMgrRequireIPReq
	var res proto.IPResMgr2IPAMRequireIPRes

	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(&req))

	res.ReqID = req.ReqID
	res.Code = proto.IPResMgrErrnoGetIPFailed
	defer sendResponse(c, &res)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {

		k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

		k8sAddrInfo := storeMgr.BindAddrInfoWithK8SResourceID(k8sResourceID, proto.K8SApiResourceKindDeployment, req.K8SPodID)
		if k8sAddrInfo == nil {
			calm_utils.Errorf("ReqID:%s get k8sAddrInfo by %s failed", req.ReqID, k8sResourceID)
		} else {
			res.IP = k8sAddrInfo.IP
			res.MacAddr = k8sAddrInfo.MacAddr
			res.SubnetGatewayAddr = k8sAddrInfo.SubNetGatewayAddr
			res.PortID = k8sAddrInfo.PortID
			res.Code = proto.IPResMgrErrnoSuccessed
			calm_utils.Debugf("ReqID:%s k8sResourceID:%d podID:%s bind with addrInfo:%s successed.", req.ReqID,
				k8sResourceID, req.K8SPodID, litter.Sdump(k8sAddrInfo))
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else {
		// 直接从nsp获取地址
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))

	return
}

func cniReleaseIP(c *gin.Context) {
	var req proto.IPAM2IPResMgrReleaseIPReq
	var res proto.IPResMgr2IPAMReleaseIPRes

	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(&req))

	res.ReqID = req.ReqID
	res.Code = proto.IPResMgrErrnoGetIPFailed
	defer sendResponse(c, &res)

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {
		err = storeMgr.UnBindAddrInfoWithK8SResourceID(k8sResourceID, proto.K8SApiResourceKindDeployment, req.K8SPodID)
		if err != nil {
			calm_utils.Errorf("ReqID:%s k8sResourceID:%d podID:%s unBind failed.", req.ReqID, k8sResourceID, req.K8SPodID)
		} else {
			calm_utils.Debugf("ReqID:%s k8sResourceID:%d podID:%s unBind successed.", req.ReqID, k8sResourceID, req.K8SPodID)
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else {
		// 处理job，cronjob
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))
	return
}
