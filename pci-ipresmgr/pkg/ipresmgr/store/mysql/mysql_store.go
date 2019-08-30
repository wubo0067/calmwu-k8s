/*
 * @Author: calm.wu
 * @Date: 2019-08-30 10:41:36
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-30 17:21:32
 */

package mysql

import (
	"context"
	"fmt"
	"pci-ipresmgr/pkg/ipresmgr/store"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var _ store.StoreMgr = &mysqlStoreMgr{}

type dbProcessHandler func(ctx context.Context) error

type mysqlStoreMgr struct {
	opts         store.StoreOptions
	dbMgr        *sqlx.DB
	mysqlConnStr string
}

func (msm *mysqlStoreMgr) doDBKeepAlive(ctx context.Context) {
	calm_utils.Debug("Start doDBKeepAlive")

	ticker := time.NewTicker(time.Minute)
	go func() {
		defer ticker.Stop()
	L:
		for {
			select {
			case <-ticker.C:
				// 定时ping
				err := msm.dbMgr.Ping()
				if err != nil {
					calm_utils.Warnf("%s connect failed. err:%s", msm.mysqlConnStr, err.Error())
				}
			case <-ctx.Done():
				calm_utils.Info("doDBKeepAlive exit")
				break L
			}
		}
		return
	}()
}

// Init 初始化
func (msm *mysqlStoreMgr) Start(ctx context.Context, opt store.Option) error {
	// 初始化参数
	opt(&msm.opts)

	// 创建mysql连接参数
	msm.mysqlConnStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", msm.opts.User, msm.opts.Passwd, msm.opts.Addr, msm.opts.DBName)

	calm_utils.Debugf("mysqlStoreMgr opts:%+v, mysqlConnStr:%s", msm.opts, msm.mysqlConnStr)

	// 创建db对象
	var err error
	msm.dbMgr, err = sqlx.Open("mysql", msm.mysqlConnStr)
	if err != nil {
		err = errors.Wrapf(err, "open %s failed.", msm.mysqlConnStr)
		calm_utils.Error(err)
		return err
	}

	// 设置默认参数
	msm.dbMgr.SetMaxIdleConns(msm.opts.IdelConnectCount)
	msm.dbMgr.SetMaxOpenConns(msm.opts.MaxOpenConnectCount)
	liftTime, err := time.ParseDuration(msm.opts.ConnectMaxLifeTime)
	if err != nil {
		err = errors.Wrapf(err, "time parse ConnectMaxLifeTime[%s] failed.", msm.opts.ConnectMaxLifeTime)
		calm_utils.Error(err.Error())
		return err
	}
	msm.dbMgr.SetConnMaxLifetime(liftTime)

	// 判断连接是否成功
	err = msm.dbMgr.Ping()
	if err != nil {
		err = errors.Wrapf(err, "%s connect failed.", msm.mysqlConnStr)
		calm_utils.Error(err.Error())
		return err
	}
	msm.doDBKeepAlive(ctx)

	calm_utils.Infof("%s connect successed", msm.mysqlConnStr)
	return nil
}

func (msm *mysqlStoreMgr) Stop() {
	if msm.dbMgr != nil {
		msm.dbMgr.Close()
	}
	return
}

func (msm *mysqlStoreMgr) dbSafeExec(ctx context.Context, dbHandler dbProcessHandler) error {
	return nil
}

func (msm *mysqlStoreMgr) RegisterSelf(instID int, listenAddr string, listenPort int) error {
	affectRows, err := msm.dbMgr.Exec(`INSERT INTO tbl_IPResMgrSrvRegister (srv_instance_name, srv_addr, register_time) VALUES (?, ?, ?)`,
		fmt.Sprintf("ipresmgr-svr-%d", instID),
		fmt.Sprintf("%s:%d", listenAddr, listenPort),
		time.Now(),
	)
	if err != nil {
		err = errors.Wrap(err, "INSERT INTO tbl_IPResMgrSrvRegister failed")
		calm_utils.Error(err)
		return err
	}

	pci
	return nil
}

// NewMysqlStoreMgr 构造一个存储对象
func NewMysqlStoreMgr() store.StoreMgr {
	mysqlStoreMgr := new(mysqlStoreMgr)
	return mysqlStoreMgr
}
