/*
 * @Author: calm.wu
 * @Date: 2019-08-28 15:25:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 16:58:46
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

	httpCode := http.StatusBadRequest
	defer func(status *int) {
		sendResponse(c, *status, &res)
	}(&httpCode)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {
		podUniqueName := makePodUniqueName(req.K8SClusterID, req.K8SNamespace, req.K8SPodName)

		k8sPodAddrInfo, err := storeMgr.BindAddrInfoWithK8SPodUniqueName(k8sResourceID, proto.K8SApiResourceKindDeployment, podUniqueName)
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
			calm_utils.Debugf("ReqID:%s k8sResourceID:%s podName:%s bind with addrInfo:%s successed.", req.ReqID,
				k8sResourceID, req.K8SPodName, litter.Sdump(k8sPodAddrInfo))

			httpCode = http.StatusCreated
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindCronJob ||
		req.K8SApiResourceKind == proto.K8SApiResourceKindJob {
		// 查询网络信息
		netRegionID, subNetID, subNetGatewayAddr, subNetCIDR, err := storeMgr.GetJobNetInfo(k8sResourceID)
		if err != nil {
			err = errors.Wrapf(err, "ReqID:%s k8sResourceID:%s GetJobNetInfo failed.", req.ReqID, k8sResourceID)
			calm_utils.Error(err.Error())
			res.Msg = err.Error()
		} else {
			// 直接从nsp获取地址
			netMask, err := getsubNetMask(subNetCIDR)
			if err != nil {
				err = errors.Wrapf(err, "ReqID:%s k8sResourceID:%s getsubNetMask failed.", req.ReqID, k8sResourceID)
				calm_utils.Error(err.Error())
				res.Msg = err.Error()
			} else {
				k8sAddrs, err := nsp.NSPMgr.AllocAddrResources(k8sResourceID, 1, netRegionID, subNetID, subNetGatewayAddr, netMask)
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
						podUniqueName := makePodUniqueName(req.K8SClusterID, req.K8SNamespace, req.K8SPodName)

						err = storeMgr.BindJobPodWithPortID(k8sResourceID, req.K8SApiResourceKind, k8sPodAddrInfo.IP, k8sPodAddrInfo.PortID, podUniqueName)
						if err != nil {
							// 归还地址
							nsp.NSPMgr.ReleaseAddrResources(k8sPodAddrInfo.PortID)
						} else {
							res.IP = k8sPodAddrInfo.IP
							res.MacAddr = k8sPodAddrInfo.MacAddr
							res.PortID = k8sPodAddrInfo.PortID
							res.SubnetGatewayAddr = k8sPodAddrInfo.SubNetGatewayAddr
							res.Code = proto.IPResMgrErrnoSuccessed
							calm_utils.Debugf("ReqID:%s k8sResourceID:%d podName:%s bind with addrInfo:%s successed.", req.ReqID,
								k8sResourceID, req.K8SPodName, litter.Sdump(k8sPodAddrInfo))
							httpCode = http.StatusCreated
						}
					}
				}
			}
		}
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))

	return
}

//  不明白，为什么bridge这个cni会重试4次，而且明确说明了不要返回error
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

	httpCode := http.StatusCreated
	defer func(status *int) {
		sendResponse(c, *status, &res)
	}(&httpCode)

	podUniqueName := makePodUniqueName(req.K8SClusterID, req.K8SNamespace, req.K8SPodName)
	// 查询pod对应的类型
	k8sResourceType := storeMgr.QueryK8SResourceKindByPodUniqueName(podUniqueName)

	if k8sResourceType == proto.K8SApiResourceKindDeployment {

		err = storeMgr.UnbindAddrInfoWithK8SPodID(proto.K8SApiResourceKindDeployment, podUniqueName)
		if err != nil {
			// TODO 告警
			err = errors.Wrapf(err, "ReqID:%s podName:%s unBind failed.", req.ReqID, req.K8SPodName)
			calm_utils.Errorf(err.Error())
			//res.Code = proto.IPResMgrErrnoReleaseIPFailed
			res.Msg = err.Error()
			//httpCode = http.StatusBadRequest
		} else {
			calm_utils.Debugf("ReqID:%s podName:%s unBind successed.", req.ReqID, req.K8SPodName)
		}
	} else if k8sResourceType == proto.K8SApiResourceKindStatefulSet {
		//
		calm_utils.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
	} else if k8sResourceType == proto.K8SApiResourceKindCronJob ||
		k8sResourceType == proto.K8SApiResourceKindJob {
		// 处理job，cronjob
		podUniqueName := makePodUniqueName(req.K8SClusterID, req.K8SNamespace, req.K8SPodName)
		err = storeMgr.UnbindJobPodWithPortID(podUniqueName)
		calm_utils.Errorf(err.Error())
		//res.Code = proto.IPResMgrErrnoReleaseIPFailed
		res.Msg = err.Error()
		//httpCode = http.StatusBadRequest
		// TODO 告警
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))
	return
}
