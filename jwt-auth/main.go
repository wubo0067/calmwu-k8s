/*
 * @Author: calm.wu
 * @Date: 2019-12-12 15:03:54
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-12-12 15:13:09
 */

// http://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac

package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	"github.com/segmentio/ksuid"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

const (
	KeyLeasePeriod = time.Duration(time.Hour * 24 * 7)
	SecretKey      = "1234567890"
)

func buildJWTToken(userName string, userPwd string) string {
	jwtPayload := jwt.MapClaims{
		"iss": "SCI",                          // token签发人
		"exp": time.Now().Add(KeyLeasePeriod), // 过期时间
		"aud": userName,                       // 受众者
		"nbf": time.Now(),                     // 签发时间
		"sub": "Deployment Helm",              // 主题
		"jti": ksuid.New().String(),           // 编号
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)
	tokenStr, err := jwtToken.SignedString(calm_utils.String2Bytes(SecretKey))
	if err != nil {
		calm_utils.Errorf("jwt SignedString failed. err:%s", err.Error())
		return ""
	}
	return tokenStr
}

func parseJWTToken(jwtToken string) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return calm_utils.String2Bytes(SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		calm_utils.Debugf("claims:%#v", litter.Sdump(claims))
		// 认为没有过期，应该返回true
		bExpired := claims.VerifyExpiresAt(time.Now().Unix(), false)
		calm_utils.Debugf("bExpired must true = [%v]", bExpired)

		// 认为过期了，应该返回false
		bExpired = claims.VerifyExpiresAt(time.Now().Unix(), true)
		calm_utils.Debugf("bExpired must false = [%v]", bExpired)
	} else {
		calm_utils.Error(err.Error())
	}
}

func ginRun() {

	go func() {
		r := gin.Default()

		r.POST("/openapi/v1/deployhelm/specifyparameters", func(c *gin.Context) {
			if jwtToken, exists := c.Request.Header["Authorization"]; exists {
				calm_utils.Debugf("jwtToken:%s", jwtToken)
				c.Status(http.StatusOK)
				return
			}
			calm_utils.Error("jwtToken not found.")
			c.Status(http.StatusUnauthorized)
		})

		r.Run("127.0.0.1:9090")
	}()
}

func sendReqAuth(jwtToken string) {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:9090/openapi/v1/deployhelm/specifyparameters", nil)
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		calm_utils.Error(err.Error())
		return
	}
	calm_utils.Debug(resp)
}

func main() {
	ginRun()
	token := buildJWTToken("ShengBin", "123456789")
	calm_utils.Debug(token)
	parseJWTToken(token)
	sendReqAuth(token)
	time.Sleep(3 * time.Second)
}
