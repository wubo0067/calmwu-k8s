package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/sanity-io/litter"
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

func testdef() {
	var i = 0
	nums := make([]int, 0)
	defer func(v *int) {
		fmt.Println("result: ", *v)
		fmt.Printf("nums: %v\n", nums)
	}(&i)

	i++
	nums = append(nums, 1)
	nums = append(nums, 2)
	nums = append(nums, 3)
}

func testGetNetMask() int {
	cidr := "192.168.1.19/26"

	pos := strings.LastIndexByte(cidr, '/')
	if pos == -1 {
		fmt.Println("/ is not found")
		return -1
	}

	fmt.Printf("/ last pos:%d, cidr len:%d\n", pos, len(cidr))

	mask, err := strconv.Atoi(cidr[pos+1:])
	if err != nil {
		fmt.Println(err.Error())
		return -2
	}
	return mask
}

func testGetClusterID() {
	k8sResourceID := "cluster-1:default:kata-nginx-deployment"
	pos := strings.IndexByte(k8sResourceID, ':')
	clusterID := k8sResourceID[:pos]
	fmt.Printf("clusterID Type-Name:%s Type-String:%s\n", reflect.ValueOf(clusterID).Type().Name(), reflect.ValueOf(clusterID).Type().String())
	fmt.Printf("clusterID :%s remaind:%s\n", clusterID, k8sResourceID[pos+1:])

	content := k8sResourceID[pos+1:]

	pos = strings.IndexByte(content, ':')
	namespace := content[:pos]
	fmt.Printf("namespace :%s\n", namespace)
}

func testDefer(flag int) {
	fmt.Printf("before defer\n")

	if flag == 1 {
		return
	}
	defer func() {
		fmt.Printf("defer after\n")
	}()
	return
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

	dumpStr := litter.Sdump(&storeCfgData)
	fmt.Println(dumpStr)

	testdef()

	testGetClusterID()

	fmt.Println(testGetNetMask())

	testDefer(1)
	testDefer(0)
}

// CGO_ENABLED=1 go build  -ldflags="-X 'main.time=`date`' -X main.version=1.0.2 -linkmode external -extldflags -static" test_v.go
