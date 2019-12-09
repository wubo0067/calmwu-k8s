/*
 * @Author: calm.wu
 * @Date: 2019-10-04 12:14:34
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 12:16:01
 */

package srv

import (
	"net/http"
	proto "pci-ipresmgr/api/proto_json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// 解绑ip和pod，job和cronjob会直接归还ip给nsp
func maintainForceUnbindIP(c *gin.Context) {
	var req proto.Maintain2IPResMgrForceUnbindIPReq
	var res proto.IPResMgr2MaintainRes

	// 解包
	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.ReqType = proto.Maintain2IPResMgrForceUnbindIPReqType
	res.Code = proto.IPResMgrErrnoSuccessed

	httpCode := http.StatusOK
	defer func(status *int, resMgr *proto.IPResMgr2MaintainRes) {
		sendResponse(c, *status, resMgr)
	}(&httpCode, &res)

	podUniqueName := makePodUniqueName(req.K8SClusterID, req.K8SNamespace, req.K8SPodName)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {
		err = storeMgr.UnbindAddrInfoWithK8SPodID(proto.K8SApiResourceKindDeployment, podUniqueName)
		if err != nil {
			// 解绑Deployment pod
			err = errors.Wrapf(err, "ReqID:%s Deployment podName:%s unBind failed.", req.ReqID, req.K8SPodName)
			calm_utils.Errorf(err.Error())
			res.Code = proto.IPResMgrErrnoMaintainForceUnbindIPFailed
			res.Msg = err.Error()
		} else {
			calm_utils.Debugf("ReqID:%s Deployment podName:%s unBind successed.", req.ReqID, req.K8SPodName)
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {
		// 解绑Statefulset pod
		res.Code = proto.IPResMgrErrnoMaintainForceUnbindIPFailed
		err = errors.Errorf("ReqID:%s not support K8SApiResourceKindStatefulSet", req.ReqID)
		calm_utils.Error(err.Error())
		res.Msg = err.Error()
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindJob ||
		req.K8SApiResourceKind == proto.K8SApiResourceKindCronJob {
		// 解绑Job Cronjob pod，这个会直接释放ip给nsp
		err = storeMgr.UnbindJobPodWithPodUniqueName(podUniqueName)
		if err != nil {
			err = errors.Wrapf(err, "ReqID:%s %s podName:%s unBind failed", req.ReqID, req.K8SApiResourceKind.String(), req.K8SPodName)
			calm_utils.Error(err.Error())
			res.Msg = err.Error()
			res.Code = proto.IPResMgrErrnoMaintainForceUnbindIPFailed
		} else {
			calm_utils.Debugf("ReqID:%s %s podName:%s unBind successed", req.ReqID, req.K8SApiResourceKind.String(), req.K8SPodName)
		}
	}

	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))
}

func maintainForceReleaseK8SResourceIPPool(c *gin.Context) {
	var req proto.Maintain2IPResMgrForceReleaseK8SResourceIPPoolReq
	var res proto.IPResMgr2MaintainRes

	// 解包
	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.ReqType = proto.Maintain2IPResMgrForceReleaseK8SResourceIPPoolReqType
	res.Code = proto.IPResMgrErrnoSuccessed

	httpCode := http.StatusOK
	defer func(status *int, resMgr *proto.IPResMgr2MaintainRes) {
		sendResponse(c, *status, resMgr)
	}(&httpCode, &res)

	k8sResourceID := makeK8SResourceID(req.K8SClusterID, req.K8SNamespace, req.K8SApiResourceName)

	if req.K8SApiResourceKind == proto.K8SApiResourceKindDeployment {
		err = storeMgr.MaintainForceReleaseK8SResourceIPPool(k8sResourceID, proto.K8SApiResourceKindDeployment)
		if err != nil {
			err = errors.Wrapf(err, "Force Release K8SResource:%s Type:%s IPPool failed.", k8sResourceID, req.K8SApiResourceKind.String())
			res.Msg = err.Error()
			calm_utils.Error(err.Error())
			res.Code = proto.IPResMgrErrnoMaintainForceReleaseK8SResourceIPPoolFailed
		} else {
			calm_utils.Debugf("Force Release K8SResource:%s Type:%s IPPool successed.", k8sResourceID, req.K8SApiResourceKind.String())
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindStatefulSet {

	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindCronJob {
		err = storeMgr.MaintainDelCronjobNetInfos(k8sResourceID)
		if err != nil {
			calm_utils.Errorf("Req:%s MaintainDelCronjobNetInfos %s failed. err:%s", req.ReqID, k8sResourceID, err.Error())
		} else {
			calm_utils.Debugf("Req:%s MaintainDelCronjobNetInfos %s successed.", req.ReqID, k8sResourceID)
		}
	} else if req.K8SApiResourceKind == proto.K8SApiResourceKindJob {
		err = storeMgr.DelJobNetInfo(k8sResourceID)
		if err != nil {
			calm_utils.Errorf("Req:%s DelJobNetInfo %s failed. err:%s", req.ReqID, k8sResourceID, err.Error())
		} else {
			calm_utils.Debugf("Req:%s DelJobNetInfo %s successed.", req.ReqID, k8sResourceID)
		}
	}
	calm_utils.Debugf("ReqID:%s Res:%s", req.ReqID, litter.Sdump(&res))
}

func maintainForceReleasePodIP(c *gin.Context) {
	var req proto.Maintain2IPResMgrForceReleasePodIPReq
	var res proto.IPResMgr2MaintainRes

	// 解包
	err := unpackRequest(c, &req)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", calm_utils.String2Bytes(err.Error()))
		return
	}
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.ReqType = proto.Maintain2IPResMgrForceReleaseK8SPodIPReqType
	res.Code = proto.IPResMgrErrnoSuccessed

	httpCode := http.StatusOK
	defer func(status *int, resMgr *proto.IPResMgr2MaintainRes) {
		sendResponse(c, *status, resMgr)
	}(&httpCode, &res)
}
