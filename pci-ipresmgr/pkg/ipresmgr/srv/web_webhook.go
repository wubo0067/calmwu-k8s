/*
 * @Author: calm.wu
 * @Date: 2019-08-28 16:03:02
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:59:16
 */

package srv

import (
	"fmt"
	"net/http"

	proto "pci-ipresmgr/api/proto_json"
	"pci-ipresmgr/pkg/ipresmgr/nsp"

	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func wbCreateIPPool(c *gin.Context) {
	var req proto.WB2IPResMgrCreateIPPoolReq
	var res proto.IPResMgr2WBRes

	// 解包
	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	// 设置返回
	res.ReqID = req.ReqID
	res.ReqType = proto.WB2IPResMgrRequestCreateIPPool
	res.Code = proto.IPResMgrErrnoCreateIPPoolFailed
	defer sendResponse(c, &res)

	if req.K8SApiResourceReplicas <= 0 {
		errInfo := fmt.Sprintf("ReqID:%s K8SApiResourceReplicas is invalid.", req.K8SApiResourceReplicas)
		calm_utils.Error(errInfo)
		res.Msg = errInfo
		return
	}

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

	// 判断是否在租期内，以及当前副本数量
	exists, replicas, err := storeMgr.CheckRecycledResources(k8sResourceID)
	if err != nil {
		res.Msg = err.Error()
		return
	}

	calm_utils.Debugf("ReqID:%s k8sResourceID:%s exists:%v replicas:%d", req.ReqID, k8sResourceID, exists, replicas)

	if !exists {
		// 从nsp获取地址
		k8sAddrs, err := nsp.NSPMgr.AllocAddrResources(k8sResourceID, &req)
		if err != nil {
			res.Msg = err.Error()
			return
		}

		// 设置地址
		err = storeMgr.SetAddrInfosToK8SResourceID(k8sResourceID, req.K8SApiResourceKind, k8sAddrs)
		if err != nil {
			// 设置对应关系失败，见IP归还给NSP
			for _, k8sAddr := range k8sAddrs {
				nsp.NSPMgr.ReleaseAddrResources(k8sAddr.PortID)
			}
			res.Msg = err.Error()
			return
		}

		res.Code = proto.IPResMgrErrnoSuccessed
		calm_utils.Infof("ReqID:%s set Addrs to k8sResourceID:%s successed.", req.ReqID, k8sResourceID)
		return
	} else {
		// 恢复的数据
		// 判断副本数量是否一致
		if req.K8SApiResourceReplicas > replicas {

		} else if req.K8SApiResourceReplicas < replicas {

		} else {
			// 副本数相同，直接返回
			res.Code = proto.IPResMgrErrnoSuccessed
		}
	}
	return
}

func wbReleaseIPPool(c *gin.Context) {
	var req proto.WB2IPResMgrReleaseIPPoolReq
	var res proto.IPResMgr2WBRes

	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.ReqType = proto.WB2IPResMgrRequestReleaseIPPool
	res.Code = proto.IPResMgrErrnoSuccessed
	defer sendResponse(c, &res)

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)
	storeMgr.AddK8SResourceAddressToRecycle(k8sResourceID, req.K8SApiResourceKind)
	return
}

func wbScaleIPPool(c *gin.Context) {
	calm_utils.Debug("wbScaleIPPool")
	c.Status(http.StatusOK)
}
