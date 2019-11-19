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
	"github.com/sanity-io/litter"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

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
