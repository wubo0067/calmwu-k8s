package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

var time string
var version string

type StoreCfgData struct {
	MysqlAddr string `json:"mysqladdr" mapstructure:"mysqladdr"`
	User      string `json:"user" mapstructure:"user"`
	Passwd    string `json:"passwd" mapstructure:"passwd"`
	DBName    string `json:"dbname" mapstructure:"dbname"`
}

func GetStoreCfgData() (StoreCfgData, string) {
	var storeCfgData StoreCfgData
	str := "hello world"
	strP := (*reflect.StringHeader)(unsafe.Pointer(&str))
	fmt.Printf("storeCfgData addr:%p str:%p strP:%+v\n", &storeCfgData, &str, strP)
	return storeCfgData, str
}

func main() {
	fmt.Println(time)
	fmt.Println(version)
	storeCfgData, str := GetStoreCfgData()
	strP := (*reflect.StringHeader)(unsafe.Pointer(&str))
	fmt.Printf("storeCfgData addr:%p str:%p strP:%+v\n", &storeCfgData, &str, strP)

	str = "Hello world"
	var v interface{} = str
	fmt.Printf("v type is:%s,%s, kind:%s\n", reflect.TypeOf(v).String(),
		reflect.ValueOf(v).Type().Name(), reflect.TypeOf(v).Kind().String())

}

// CGO_ENABLED=1 go build  -ldflags="-X 'main.time=`date`' -X main.version=1.0.2 -linkmode external -extldflags -static" test_v.go
