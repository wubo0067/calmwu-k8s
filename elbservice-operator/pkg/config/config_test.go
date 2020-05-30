/*
 * @Author: calmwu
 * @Date: 2020-05-30 23:31:38
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-05-30 23:48:37
 */

package config

import (
	"testing"

	"github.com/sanity-io/litter"
)

func TestReadConfig(t *testing.T) {
	err := unmarshalCfgFile("../../deploy/config/elbservice_config.json")
	if err != nil {
		t.Errorf("unmarshalCfgFile failed, err:%s\n", err.Error())
	}

	t.Logf("ELBServiceConfig: %s", litter.Sdump(ELBServiceCfg))
}
