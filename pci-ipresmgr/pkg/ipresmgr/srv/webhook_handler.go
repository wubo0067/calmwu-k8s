/*
 * @Author: calm.wu
 * @Date: 2019-08-28 16:03:02
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-28 17:19:16
 */

package srv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func wbCreateIPPool(c *gin.Context) {
	calm_utils.Debug("wbCreateIPPool")
	c.Status(http.StatusOK)
}

func wbReleaseIPPool(c *gin.Context) {
	calm_utils.Debug("wbReleaseIPPool")
	c.Status(http.StatusOK)
}

func wbScaleIPPool(c *gin.Context) {
	calm_utils.Debug("wbScaleIPPool")
	c.Status(http.StatusOK)
}
