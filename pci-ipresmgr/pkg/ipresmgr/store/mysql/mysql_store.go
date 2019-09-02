/*
 * @Author: calm.wu
 * @Date: 2019-08-30 10:41:36
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-01 09:52:51
 */

package mysql

import (
	"context"
	"fmt"
	"pci-ipresmgr/pkg/ipresmgr/store"
	"pci-ipresmgr/table"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var _ store.StoreMgr = &mysqlStoreMgr{}

type dbProcessHandler func(ctx context.Context) error

type mysqlStoreMgr struct {
	opts                       store.StoreOptions
	dbMgr                      *sqlx.DB
	mysqlConnStr               string
	addrResourceLeasePeriodMgr AddrResourceLeasePeriodMgrItf
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

	msm.addrResourceLeasePeriodMgr = NewAddrResourceLeasePeriodMgr(ctx, msm.dbMgr)

	calm_utils.Infof("%s connect successed", msm.mysqlConnStr)
	return nil
}

func (msm *mysqlStoreMgr) Stop() {
	if msm.dbMgr != nil {
		msm.dbMgr.Close()
		msm.dbMgr = nil
		calm_utils.Info("mysqlStoreMgr stop")
	}
	return
}

func (msm *mysqlStoreMgr) dbSafeExec(ctx context.Context, dbHandler dbProcessHandler) (dbExecErr error) {
	defer func() {
		if err := recover(); err != nil {
			stackInfo := calm_utils.CallStack(3)
			dbExecErr = errors.Errorf("Panic! err:%v stack:%s", err, stackInfo)
			calm_utils.Error(err)
		}
	}()

	return dbHandler(ctx)
}

func (msm *mysqlStoreMgr) Register(listenAddr string, listenPort int) error {
	vCtx := setCtxVal(context.Background(), "listenAddr", listenAddr)
	vCtx = setCtxVal(vCtx, "listenPort", listenPort)

	return msm.dbSafeExec(vCtx,
		func(ctx context.Context) error {
			listenAddr, _ := getCtxStrVal(ctx, "listenAddr")
			listenPort, _ := getCtxIntVal(ctx, "listenPort")
			srvAddr := fmt.Sprintf("%s:%d", listenAddr, listenPort)
			registerTime := time.Now()

			_, err := msm.dbMgr.Exec("INSERT INTO tbl_IPResMgrSrvRegister (srv_instance_name, srv_addr, register_time) VALUES (?, ?, ?)",
				msm.opts.SrvInstID, srvAddr, registerTime)
			if err != nil {
				err = errors.Wrap(err, "INSERT INTO tbl_IPResMgrSrvRegister failed")
				calm_utils.Error(err)
				return err
			}
			calm_utils.Infof("Register %s successed.", msm.opts.SrvInstID)

			// 开始加载tbl_K8SResourceIPRecycle
			addrRows, err := msm.dbMgr.Queryx("SELECT * FROM tbl_K8SResourceIPRecycle WHERE srv_instance_name=?", msm.opts.SrvInstID)
			if err != nil {
				err = errors.Wrapf(err, "SELECT * FROM tbl_K8SResourceIPRecycle WHERE srv_instance_name=%s failed.", msm.opts.SrvInstID)
				calm_utils.Error(err)
				return err
			}

			loadCount := 0
			for addrRows.Next() {
				addrRecyclingRecord := new(table.TblK8SResourceIPRecycleS)
				err = addrRows.StructScan(addrRecyclingRecord)
				if err != nil {
					calm_utils.Fatalf("Scan TblK8SResourceIPRecycleS failed. err:%s", err.Error())
				}
				loadCount++
				msm.addrResourceLeasePeriodMgr.AddLeaseRecyclingRecord(addrRecyclingRecord)
			}
			calm_utils.Infof("load from tbl_K8SResourceIPRecycle %d records", loadCount)
			// 加载完毕，启动
			err = msm.addrResourceLeasePeriodMgr.Start()
			if err != nil {
				calm_utils.Fatalf("mysqlAddrResourceLeasePeriodMgr start failed. err:%s", err.Error())
			}
			calm_utils.Info("mysqlAddrResourceLeasePeriodMgr start successed.")
			return nil
		},
	)
}

func (msm *mysqlStoreMgr) UnRegister() {
	msm.dbSafeExec(context.Background(),
		func(ctx context.Context) error {
			_, err := msm.dbMgr.Exec("DELETE FROM tbl_IPResMgrSrvRegister WHERE srv_instance_name=?",
				msm.opts.SrvInstID)
			if err != nil {
				err = errors.Wrapf(err, "DELETE FROM tbl_IPResMgrSrvRegister WHERE srv_instance_name='%s' failed.", msm.opts.SrvInstID)
				calm_utils.Error(err)
				return err
			}
			calm_utils.Infof("unRegister %s successed.", msm.opts.SrvInstID)
			return nil
		},
	)
	return
}

// NewMysqlStoreMgr 构造一个存储对象
func NewMysqlStoreMgr() store.StoreMgr {
	mysqlStoreMgr := new(mysqlStoreMgr)
	return mysqlStoreMgr
}
