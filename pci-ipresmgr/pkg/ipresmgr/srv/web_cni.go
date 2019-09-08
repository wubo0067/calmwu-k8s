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
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.Code = proto.IPResMgrErrnoGetIPFailed
	defer sendResponse(c, &res)

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
	calm_utils.Debugf("Req:%s", litter.Sdump(req))

	res.ReqID = req.ReqID
	res.Code = proto.IPResMgrErrnoGetIPFailed
	defer sendResponse(c, &res)

	return
}
