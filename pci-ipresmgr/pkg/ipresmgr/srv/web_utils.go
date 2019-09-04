/*
 * @Author: calm.wu
 * @Date: 2019-09-04 14:20:42
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 15:03:59
 */

package srv

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func unpackRequest(c *gin.Context, req interface{}) error {
	bodyData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		err = errors.Wrap(err, "unpackRequest failed.")
		calm_utils.Error(err)
		return err
	}

	calm_utils.Debugf("request:%s", calm_utils.Bytes2String(bodyData))

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(bodyData, req)
	if err != nil {
		err = errors.Wrap(err, "json Unmarshal failed.")
		calm_utils.Error(err)
		return err
	}
	return nil
}

func sendResponse(c *gin.Context, res interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	resData, err := json.Marshal(res)
	if err != nil {
		calm_utils.Error(err)
		return
	}
	calm_utils.Debugf("send response to %s successed", c.Request.RemoteAddr)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", resData)
	return
}
