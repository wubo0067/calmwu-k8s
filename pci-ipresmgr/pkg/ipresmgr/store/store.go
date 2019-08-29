/*
 * @Author: calm.wu
 * @Date: 2019-08-29 18:44:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-29 19:21:57
 */

package store

import (
	"github.com/jmoiron/sqlx"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// Store 存储接口
type Store interface {
	Start(Option) error
	Stop()
}

// StoreOptions 存储的参数
type StoreOptions struct {
	Addr   string
	User   string
	Passwd string
	DBName string
}

// Option 选项修改
type Option func(*StoreOptions)

type backendStore struct {
	opts  StoreOptions
	dbMgr *sqlx.DB
}

// Init 初始化
func (bs *backendStore) Start(opt Option) error {
	opt(&bs.opts)

	calm_utils.Debugf("backendStore opts:%+v", bs.opts)

	return nil
}

func (bs *backendStore) Stop() {
	return
}

// NewStore 构造一个存储对象
func NewStore() Store {
	return new(backendStore)
}
