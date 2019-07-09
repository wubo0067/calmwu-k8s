/*
 * @Author: calm.wu
 * @Date: 2019-07-04 11:16:04
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-09 11:20:43
 */

package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewZapLog(t *testing.T) {
	logFullName := "./test.log"
	logLevel := zapcore.DebugLevel

	log := NewZapLog(logFullName, logLevel)
	if log == nil {
		t.Error("NewZapLog invoke failed!")
		return
	}

	i := 0
	for i < 10 {
		log.Debugf("Debugf %d", i)
		log.Infof("Infof %d", i)
		log.Warnf("Infof %d", i)
		i++
	}
}
