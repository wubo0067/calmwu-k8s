/*
 * @Author: CALM.WU
 * @Date: 2021-04-14 10:29:22
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2021-04-14 17:57:58
 */

package main

import (
	"os"

	"dynifr-watchres/src/config"
	"dynifr-watchres/src/kubehelper"
	"dynifr-watchres/src/kubewatch"

	"github.com/sanity-io/litter"
	calmUtils "github.com/wubo0067/calmwu-go/utils"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	if err := config.LoadConf(); err != nil {
		calmUtils.Error(err.Error())
		os.Exit(-1)
	}

	// 信号管理
	stopCtx := signals.SetupSignalHandler()

	// 获取配置数据
	confData := config.GetConfData()
	calmUtils.Debugf("config info: %s", litter.Sdump(confData))

	// create dynamic client
	dc := kubehelper.MakeDynamicClient(confData.KubeCfg)
	calmUtils.Debugf("cluster dc: %s", litter.Sdump(dc))

	kubewatch.DynamicInformerWatchResources(dc, stopCtx.Done())

	<-stopCtx.Done()
	calmUtils.Debug("dynifr-watchres exit!")
}
