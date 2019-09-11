/*
 * @Author: calm.wu
 * @Date: 2019-09-07 16:07:14
 * @Last Modified by:   calm.wu
 * @Last Modified time: 2019-09-07 16:07:14
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	proto "pci-ipresmgr/api/proto_json"

	"github.com/dghubble/sling"
	jsoniter "github.com/json-iterator/go"
	"github.com/sanity-io/litter"
	"github.com/segmentio/ksuid"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	srvIPResMgrAddr = flag.String("svraddr", "http://192.168.6.134:30001/", "srv ipresmgr addr")
	testType        = flag.Int("type", 1, "1: CreateIPPool, 2: ReleaseIPPool, 3: ScaleIPPool, 4: RequireIP, 5: ReleaseIP")
	logger          *log.Logger
)

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func testCreateIPPool() {
	var createIPPoolReq proto.WB2IPResMgrCreateIPPoolReq
	createIPPoolReq.ReqID = ksuid.New().String()
	createIPPoolReq.K8SApiResourceKind = proto.K8SApiResourceKindDeployment
	createIPPoolReq.K8SClusterID = "cluster-1"
	createIPPoolReq.K8SNamespace = "default"
	createIPPoolReq.K8SApiResourceName = "kata-nginx-deployment"
	createIPPoolReq.K8SApiResourceReplicas = 3
	createIPPoolReq.NetRegionalID = fmt.Sprintf("netregional-%s", ksuid.New().String())
	createIPPoolReq.SubnetID = fmt.Sprintf("subnet-%s", ksuid.New().String())
	createIPPoolReq.SubnetGatewayAddr = "1.1.1.1"

	var createIPPooRes proto.IPResMgr2WBRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&createIPPoolReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ippool/").Post("create").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&createIPPooRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ippool/create failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode != 200 {
		logger.Fatalf("post %sv1/ippool/create failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("createIPPooRes:%s\n", litter.Sdump(&createIPPooRes))
}

func testReleaseIPPool() {
	var releaseIPPoolReq proto.WB2IPResMgrReleaseIPPoolReq
	releaseIPPoolReq.ReqID = ksuid.New().String()
	releaseIPPoolReq.K8SApiResourceKind = proto.K8SApiResourceKindDeployment
	releaseIPPoolReq.K8SClusterID = "cluster-1"
	releaseIPPoolReq.K8SApiResourceName = "default"
	releaseIPPoolReq.K8SApiResourceName = "kata-nginx-deployment"

	var releaseIPPooRes proto.IPResMgr2WBRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&releaseIPPoolReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ippool/").Post("release").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&releaseIPPooRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ippool/release failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode != 200 {
		logger.Fatalf("post %sv1/ippool/release failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("releaseIPPooRes:%s\n", litter.Sdump(&releaseIPPooRes))
}

func main() {
	flag.Parse()

	logger = calm_utils.NewSimpleLog(nil)

	switch *testType {
	case 1:
		testCreateIPPool()
	case 2:
		testReleaseIPPool()
	default:
		logger.Fatalf("Not support type:%d\n", *testType)
	}

	logger.Println("test completed")
}
