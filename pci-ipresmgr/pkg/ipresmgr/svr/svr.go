/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:45:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-26 19:51:06
 */

package svr

import "github.com/micro/cli"

var (
	// SvrFlags 命令参数
	SvrFlags = []cli.Flag{
		cli.IntFlag{
			Name:  "id",
			Value: 1,
			Usage: "IPResMgrSvr instance ID",
		},
		cli.StringFlag{
			Name:  "logpath",
			Value: "log",
			Usage: "",
		},
	}
)

// SvrMain 服务的入口
func SvrMain(c *cli.Context) error {
	return nil
}
