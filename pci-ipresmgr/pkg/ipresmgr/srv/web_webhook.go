/*
 * @Author: calm.wu
 * @Date: 2019-08-28 16:03:02
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 15:24:37
 */

package srv

import (
	"net/http"

	proto "pci-ipresmgr/api/proto_json"

	"github.com/gin-gonic/gin"
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

	// 设置返回
	res.ReqID = req.ReqID
	res.ReqType = proto.WB2IPResMgrRequestCreateIPPool
	res.Code = proto.IPResMgrErrnoCreateIPPoolFailed
	defer sendResponse(c, &res)

	// 判断是否在租期内，以及当前副本数量
}

func wbReleaseIPPool(c *gin.Context) {
	calm_utils.Debug("wbReleaseIPPool")
	c.Status(http.StatusOK)
}

func wbScaleIPPool(c *gin.Context) {
	calm_utils.Debug("wbScaleIPPool")
	c.Status(http.StatusOK)
}
