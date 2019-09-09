/*
 * @Author: calm.wu
 * @Date: 2019-08-31 17:31:02
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-31 17:51:35
 */

package mysql

import (
	"context"

	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

type contextKey string

func setCtxVal(ctx context.Context, key string, val interface{}) context.Context {
	valCtx := context.WithValue(ctx, contextKey(key), val)
	return valCtx
}

func getCtxStrVal(ctx context.Context, key string) (string, error) {
	v, ok := ctx.Value(contextKey(key)).(string)
	if !ok {
		err := errors.Errorf("key:%s value is not string", key)
		calm_utils.Error(err)
		return "", err
	}
	return v, nil
}

func getCtxIntVal(ctx context.Context, key string) (int, error) {
	v, ok := ctx.Value(contextKey(key)).(int)
	if !ok {
		err := errors.Errorf("key:%s value is not int", key)
		calm_utils.Error(err)
		return 0, err
	}
	return v, nil
}
