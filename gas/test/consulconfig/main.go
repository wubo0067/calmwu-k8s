/*
 * @Author: calm.wu
 * @Date: 2019-07-03 16:53:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-03 16:55:48
 */

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gas/internal/utils/config"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func main() {
	keyPath := "/micro/config"
	fmt.Printf("path:%v\n", strings.Split(keyPath, "/"))

	cc, err := config.NewConfigMgr(keyPath)
	if err != nil {
		log.Printf("NewConfigMgr micro/config/cache failed! reason:%s", err.Error())
		return
	}

	configData := cc.GetConfigData("micro", "config", "cache")
	fmt.Printf("configData:%s\n", calm_utils.Bytes2String(configData))

	go func() {
	L:
		for {
			select {
			case _, ok := <-cc.ChangeNtfChan:
				if ok {
					fmt.Printf("receive config change notify\n")
				} else {
					fmt.Printf("change notify channel is closed!\n")
					break L
				}
			}
		}
	}()

	time.Sleep(10 * time.Second)
	configData = cc.GetConfigData("micro", "config", "database")
	fmt.Printf("configData:%s\n", calm_utils.Bytes2String(configData))

	cc.Stop()
}
