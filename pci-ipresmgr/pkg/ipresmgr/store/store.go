/*
 * @Author: calm.wu
 * @Date: 2019-08-29 18:44:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-30 17:21:23
 */

package store

import "context"

// StoreMgr 存储接口
type StoreMgr interface {
	// Start 启动存储管理
	Start(context.Context, Option) error
	// Stop 停止存储管理
	Stop()
	// 注册自己，保证实例id唯一
	RegisterSelf(instID string, listenAddr string, listenPort int) error
	// 注销自己
	UnRegisterSelf(instID string)
}

// StoreOptions 存储的参数
type StoreOptions struct {
	Addr                string
	User                string
	Passwd              string
	DBName              string
	IdelConnectCount    int
	MaxOpenConnectCount int
	ConnectMaxLifeTime  string
}

// Option 选项修改
type Option func(*StoreOptions)
