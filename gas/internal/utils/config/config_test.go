/*
 * @Author: calm.wu
 * @Date: 2019-07-03 14:21:00
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-03 16:53:53
 */

package config

import (
	"fmt"
	"testing"
	"time"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

func TestReadConsulConfig(t *testing.T) {
	keyPath := "/micro/config/cache"

	config, err := NewConfigMgr(keyPath)
	if err != nil {
		t.Errorf("NewConfigMgr micro/config/cache failed! reason:%s", err.Error())
		return
	}

	configData := config.GetConfigData()
	fmt.Printf("configData:%s\n", calm_utils.Bytes2String(configData))

	go func() {
	L:
		for {
			select {
			case _, ok := <-config.ChangeNtfChan:
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
	// 这个不在关注的keypath中，所以获取不到数据
	configData = config.GetConfigData("micro", "config", "database")
	fmt.Printf("configData:%s\n", calm_utils.Bytes2String(configData))

	config.Stop()
}
