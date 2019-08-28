/*
 * @Author: calm.wu
 * @Date: 2019-08-28 15:25:30
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-28 17:18:57
 */

package srv

import (
	"net/http"

	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func cniRequireIP(c *gin.Context) {
	calm_utils.Debug("cniRequireIP")
	c.Status(http.StatusOK)
}

func cniReleaseIP(c *gin.Context) {
	calm_utils.Debug("cniReleaseIP")
	c.Status(http.StatusOK)
}
