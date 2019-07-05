/*
 * @Author: calm.wu
 * @Date: 2019-07-04 11:01:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-07-04 11:20:20
 */

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// 创建一个ZapLog对象
func NewZapLog(logFullName string, logLevel zapcore.Level) *zap.SugaredLogger {
	calm_utils.InitDefaultZapLog(logFullName, logLevel)
	return calm_utils.ZLog
}
