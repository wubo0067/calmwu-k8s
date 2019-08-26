/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:23:59
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-26 16:45:24
 */

package main

import (
	"pci-ipresmgr/pkg/ipresmgr"
)

var (
	// Version 版本
	Version string
	// BuildTime 时间
	BuildTime string
)

func main() {
	ipresmgr.Main(Version, BuildTime)
}
