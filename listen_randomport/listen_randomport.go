/*
 * @Author: calm.wu
 * @Date: 2019-09-23 14:49:58
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-23 17:18:40
 */

package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func startGin() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		log.Println(example)
		c.Writer.WriteHeader(http.StatusOK)
	})

	// :和*的区别，action获得值有/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		log.Println(message)
		c.String(http.StatusOK, message)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run("0.0.0.0:0")
}

func main() {
	logger := calm_utils.NewSimpleLog(nil)
	listener, err := net.Listen("tcp", "0.0.0.0:0")

	if err != nil {
		logger.Fatal(err.Error())
	}

	port := listener.Addr().(*net.TCPAddr).Port
	logger.Printf("Listen on port:%d\n", port)

	startGin()
}
