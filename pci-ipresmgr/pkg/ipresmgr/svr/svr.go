/*
 * @Author: calm.wu
 * @Date: 2019-08-26 14:45:38
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-26 15:12:53
 */

package svr

import "github.com/micro/cli"

var (
	// IPResMgrSvrFlags 命令参数
	IPResMgrSvrFlags = []cli.Flag{
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

// IPResMgrSvrMain 服务的入口
func IPResMgrSvrMain(c *cli.Context) error {
	return nil
}
