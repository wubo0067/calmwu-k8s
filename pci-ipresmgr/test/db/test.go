/*
 * @Author: calm.wu
 * @Date: 2019-08-29 14:29:52
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:38:38
 */

// https://cloud.tencent.com/developer/article/1079583
// http://jmoiron.github.io/sqlx/
// https://github.com/jmoiron/sqlx

package main

import (
	"database/sql"
	"fmt"
	"log"
	"pci-ipresmgr/table"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

type TblTestS struct {
	ID                     uint           `db:"id"`
	K8SResourceID          string         `db:"k8sresource_id"`
	CreateTime             time.Time      `db:"create_time"`
	NSPResourceReleaseTime time.Time      `db:"nspresource_release_time"`
	SubNetID               sql.NullString `db:"subnet_id"`
	//SubNetID     string `db:"subnet_id"`
	NspResources []byte `db:"nsp_resources"`
}

func initDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "root:root@tcp(192.168.6.134:3306)/db_ipresmgr?parseTime=true&loc=Local")
	if err != nil {
		log.Fatalf("root:root@tcp(192.168.6.134:3306)/db_ipresmgr failed. err:\n", err.Error())
	}
	log.Println("root:root@tcp(192.168.6.134:3306)/db_ipresmgr open successed")

	// 缺省设置
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(16)
	db.SetConnMaxLifetime(time.Hour * 24 * 7)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Connect root:root@tcp(192.168.6.134:3306)/db_ipresmgr failed. err:%s\n", err.Error())
	}
	log.Println("root:root@tcp(192.168.6.134:3306)/db_ipresmgr connect successed")
	return db
}

func insertTbltest(db *sqlx.DB) {

	defer func() {
		if err := recover(); err != nil {
			stackInfo := calm_utils.CallStack(3)
			log.Printf("err:%s stack:%s", err, stackInfo)
		}
	}()

	nspResource := struct {
		IP  string
		Mac string
	}{
		IP:  "192.168.1.1",
		Mac: "00:0c:29:7a:9d:78",
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var err error

	var test TblTestS
	test.K8SResourceID = "k8sclusterid-namespace-resource_name"
	test.NspResources, err = json.Marshal(nspResource)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO tbl_Test 
	(k8sresource_id, nspresource_release_time, subnet_id, create_time, nsp_resources) VALUES 
	(?, ?, ?, ?, ?)`,
		"test-1",
		time.Now().Add(time.Hour),
		uuid.New().String(),
		time.Now(),
		test.NspResources)
	if err != nil {
		log.Fatal(err)
	}

	// 测试null值的处理
	tx := db.MustBegin()
	for i := 0; i < 3; i++ {
		sqlRes := tx.MustExec(`INSERT INTO tbl_Test (k8sresource_id, nspresource_release_time, nsp_resources) VALUES (?, ?, ?)`,
			fmt.Sprintf("%s-%d", test.K8SResourceID, i),
			time.Now().Add(time.Hour),
			test.NspResources)
		affectedRows, err := sqlRes.RowsAffected()
		if err != nil {
			log.Printf("insert failed. err:%s\n", err.Error())
			tx.Rollback()
			return
		}
		log.Printf("insert successed. affectedRows:%d\n", affectedRows)

		// 做个rollback测试，前面插入的会回滚掉
		// if i == 8 {
		// 	log.Printf("do rollback at %d\n", i)
		// 	tx.Rollback()
		// 	return
		// }
	}
	tx.Commit()
}

func selectTbltest(db *sqlx.DB) {
	var testLst []*TblTestS
	err := db.Select(&testLst, "SELECT * FROM tbl_Test")
	if err != nil {
		log.Fatalf("select failed. %s\n", err.Error())
	}

	for _, test := range testLst {
		log.Printf("test K8SResourceID:%s releaseTime:%s subnetid:[%s] createTime:%s nspResource:%s\n",
			test.K8SResourceID,
			test.NSPResourceReleaseTime.String(),
			test.SubNetID.String,
			//test.SubNetID,
			test.CreateTime.String(),
			calm_utils.Bytes2String(test.NspResources))

		log.Printf("createtime is zero:%v\n", test.CreateTime.IsZero())
	}
}

func insertMultilRecored(db *sqlx.DB) {
	nspResource := struct {
		IP  string
		Mac string
	}{
		IP:  "192.168.1.1",
		Mac: "00:0c:29:7a:9d:78",
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var err error
	var test TblTestS
	test.K8SResourceID = "k8sclusterid-namespace-resource_name"
	test.NspResources, err = json.Marshal(nspResource)

	for i := 0; i < 10; i++ {
		_, err = db.Exec(`INSERT INTO tbl_Test 
		(k8sresource_id, nspresource_release_time, subnet_id, create_time, nsp_resources) VALUES 
		(?, ?, ?, ?, ?)`,
			fmt.Sprintf("test-%d", i),
			time.Now().Add(time.Hour),
			uuid.New().String(),
			time.Now(),
			test.NspResources)
		if err != nil {
			log.Fatalf("insert tbl_Test %d failed. err:%s\n", i, err.Error())
		}
	}
	log.Printf("insert 1000 recored successed!\n")
}

func testScanRows(db *sqlx.DB) {
	rows, err := db.Queryx("SELECT * FROM tbl_Test WHERE k8sresource_id=?", "test-0")
	for rows.Next() {
		var test TblTestS
		err = rows.StructScan(&test)
		if err != nil {
			log.Fatalf("row StructScan failed. err:%s\n", err.Error())
		}
		log.Printf("test:%+v\n", test)
	}
}

func deleteInvalidRow(db *sqlx.DB) {
	sqlResult, err := db.Exec("DELETE FROM tbl_Test WHERE k8sresource_id=?", "test-0111")
	if err != nil {
		log.Fatalf("exec delete invalid test-0111 record failed. err:%s", err.Error())
	}
	affectRows, err := sqlResult.RowsAffected()
	if err != nil {
		log.Fatalf("delete invalid test-0111 RowsAffected failed. err:%s", err.Error())
	}
	// 不存在的记录只能通过affectRows来判断。这不是error
	log.Printf("exec delete invalid test-0111 affectRows:%d", affectRows)
}

func insertMultiK8SResourceIPRecycles(db *sqlx.DB) {
	log.Println(time.Now().String())

	var recycleRecord table.TblK8SResourceIPRecycleS
	recycleRecord.SrvInstanceName = "ipresmgr-svr_1"
	recycleRecord.CreateTime = time.Now()

	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO tbl_K8SResourceIPRecycle 
		(srv_instance_name, k8sresource_id, k8sresource_type, replicas, create_time, nspresource_release_time) VALUES 
		(?, ?, ?, ?, ?, ?)`,
			recycleRecord.SrvInstanceName,
			fmt.Sprintf("k8sresource-%d", i),
			i,
			i,
			recycleRecord.CreateTime,
			recycleRecord.CreateTime.Add(time.Duration(60+i*5)*time.Second),
		)
		if err != nil {
			log.Fatalf("insert tbl_K8SResourceIPRecycle %d failed. err:%s\n", i, err.Error())
		}
	}
}

func testQueryColumn(db *sqlx.DB) {
	var subnet_id string
	err := db.Get(&subnet_id, "SELECT subnet_id FROM tbl_Test WHERE k8sresource_id=?", "test-0")
	if err != nil {
		log.Fatalf("err: %s", err.Error())
	}
	log.Printf("tbl_Test.subnet_id:%s\n", subnet_id)
}

func main() {
	calm_utils.NewSimpleLog(nil)

	db := initDB()
	defer db.Close()

	//insertTbltest(db)
	//selectTbltest(db)
	insertMultilRecored(db)
	//testScanRows(db)
	//deleteInvalidRow(db)
	//insertMultiK8SResourceIPRecycles(db)
	//testQueryColumn(db)
}
