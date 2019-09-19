/*
 * @Author: calm.wu
 * @Date: 2019-08-28 15:25:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-13 10:57:48
 */

package srv

import (
	"net/http"

	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/nsp"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

	res.ReqID = req.ReqID
	res.Code = proto.IPResMgrErrnoGetIPFailed
	defer sendResponse(c, &res)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {

		k8sPodAddrInfo, err := storeMgr.BindAddrInfoWithK8SPodID(k8sResourceID, proto.K8SApiResourceKindDeployment, req.K8SPodID)
		if err != nil {
			err := errors.Wrapf(err, "ReqID:%s get k8sPodAddrInfo by %s failed", req.ReqID, k8sResourceID)
			calm_utils.Error(err.Error())
			res.Msg = err.Error()
		} else {
			res.IP = k8sPodAddrInfo.IP
			res.MacAddr = k8sPodAddrInfo.MacAddr
			res.SubnetGatewayAddr = k8sPodAddrInfo.SubNetGatewayAddr
			res.PortID = k8sPodAddrInfo.PortID
			res.Code = proto.IPResMgrErrnoSuccessed
			calm_utils.Debugf("ReqID:%s k8sResourceID:%s podID:%s bind with addrInfo:%s successed.", req.ReqID,
				k8sResourceID, req.K8SPodID, litter.Sdump(k8sPodAddrInfo))
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else {
		// 查询网络信息
		netRegionID, subNetID, subNetGatewayAddr, err := storeMgr.GetJobNetInfo(k8sResourceID)
		if err != nil {
			err = errors.Wrapf(err, "ReqID:%s GetJobNetInfo %s failed.", req.ReqID, k8sResourceID)
			calm_utils.Error(err.Error())
			res.Msg = err.Error()
		} else {
			// 直接从nsp获取地址
			k8sAddrs, err := nsp.NSPMgr.AllocAddrResources(k8sResourceID, 1, netRegionID, subNetID, subNetGatewayAddr)
			if err != nil {
				err = errors.Wrapf(err, "ReqID:%s AllocAddrResources from nsp %s failed.", req.ReqID, k8sResourceID)
				calm_utils.Error(err.Error())
				res.Msg = err.Error()
			} else {
				allocK8SPodAddrsSize := len(k8sAddrs)
				if allocK8SPodAddrsSize < 1 {
					err = errors.Wrapf(err, "ReqID:%s AllocAddrResources count:%d from nsp %s is invalid.", req.ReqID, allocK8SPodAddrsSize, k8sResourceID)
					calm_utils.Error(err.Error())
					res.Msg = err.Error()
				} else {
					k8sPodAddrInfo := k8sAddrs[0]
					err = storeMgr.BindJobPodWithPortID(k8sResourceID, k8sPodAddrInfo.IP, k8sPodAddrInfo.PortID, req.K8SPodID)
					if err != nil {
						// 归还地址
						nsp.NSPMgr.ReleaseAddrResources(k8sPodAddrInfo.PortID)
					} else {
						res.IP = k8sPodAddrInfo.IP
						res.MacAddr = k8sPodAddrInfo.MacAddr
						res.PortID = k8sPodAddrInfo.PortID
						res.SubnetGatewayAddr = k8sPodAddrInfo.SubNetGatewayAddr
						res.Code = proto.IPResMgrErrnoSuccessed
						calm_utils.Debugf("ReqID:%s k8sResourceID:%d podID:%s bind with addrInfo:%s successed.", req.ReqID,
							k8sResourceID, req.K8SPodID, litter.Sdump(k8sPodAddrInfo))
					}
				}
			}
		}
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
	res.Code = proto.IPResMgrErrnoSuccessed
	defer sendResponse(c, &res)

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {
		err = storeMgr.UnbindAddrInfoWithK8SPodID(k8sResourceID, proto.K8SApiResourceKindDeployment, req.K8SPodID)
		if err != nil {
			// TODO 告警
			calm_utils.Errorf("ReqID:%s k8sResourceID:%s podID:%s unBind failed.", req.ReqID, k8sResourceID, req.K8SPodID)
		} else {
			calm_utils.Debugf("ReqID:%s k8sResourceID:%s podID:%s unBind successed.", req.ReqID, k8sResourceID, req.K8SPodID)
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else {
		// 处理job，cronjob
		storeMgr.UnbindJobPodWithPortID(k8sResourceID, req.K8SPodID)
		// TODO 告警
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))
	return
}
