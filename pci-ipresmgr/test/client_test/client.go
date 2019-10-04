/*
 * @Author: calm.wu
 * @Date: 2019-09-07 16:07:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 16:50:06
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	proto "pci-ipresmgr/api/proto_json"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/sanity-io/litter"
	"github.com/segmentio/ksuid"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

var (
	srvIPResMgrAddr = flag.String("svraddr", "http://192.168.6.134:30001/", "srv ipresmgr addr")
	testType        = flag.Int("type", 1, "1: CreateIPPool, 2: ReleaseIPPool, 3: ScaleIPPool, 4: RequireIP, 5: ReleaseIP")
	unBindPodName   = flag.String("unbindpodname", "", "Unbind podName")
	oldReplicas     = flag.Int("oldreplicas", 0, "old replicas")
	newReplicas     = flag.Int("newreplicas", 1, "new replicas")
	parallel        = flag.Int("parallel", 1, "parallel requests")
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
	createIPPoolReq.SubnetCIDR = "1.1.1.1/26"

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
	releaseIPPoolReq.K8SNamespace = "default"
	releaseIPPoolReq.K8SApiResourceName = "kata-nginx-deployment"

	var releaseIPPooRes proto.IPResMgr2WBRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&releaseIPPoolReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ippool/").Post("release").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&releaseIPPooRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ippool/release failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		logger.Fatalf("post %sv1/ippool/release failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("releaseIPPooRes:%s\n", litter.Sdump(&releaseIPPooRes))
}

func testScaleIPPool() {
	var scaleIPPoolReq proto.WB2IPResMgrScaleIPPoolReq
	scaleIPPoolReq.ReqID = ksuid.New().String()
	scaleIPPoolReq.K8SApiResourceKind = proto.K8SApiResourceKindDeployment
	scaleIPPoolReq.K8SClusterID = "cluster-1"
	scaleIPPoolReq.K8SNamespace = "default"
	scaleIPPoolReq.K8SApiResourceName = "kata-nginx-deployment"
	scaleIPPoolReq.K8SApiResourceOldReplicas = *oldReplicas
	scaleIPPoolReq.K8SApiResourceNewReplicas = *newReplicas
	scaleIPPoolReq.NetRegionalID = fmt.Sprintf("netregional-%s", ksuid.New().String())
	scaleIPPoolReq.SubnetID = fmt.Sprintf("subnet-%s", ksuid.New().String())
	scaleIPPoolReq.SubnetGatewayAddr = "1.1.1.1"
	scaleIPPoolReq.SubnetCIDR = "1.1.1.1/26"

	var scaleIPPoolRes proto.IPResMgr2WBRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&scaleIPPoolReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ippool/").Post("scale").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&scaleIPPoolRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ippool/scale failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		logger.Fatalf("post %sv1/ippool/scale failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("scaleIPPoolRes:%s\n", litter.Sdump(&scaleIPPoolRes))
}

func testRequireIP() {
	var requireIPReq proto.IPAM2IPResMgrRequireIPReq
	requireIPReq.ReqID = ksuid.New().String()
	requireIPReq.K8SApiResourceKind = proto.K8SApiResourceKindDeployment
	requireIPReq.K8SClusterID = "cluster-1"
	requireIPReq.K8SNamespace = "default"
	requireIPReq.K8SApiResourceName = "kata-nginx-deployment"
	requireIPReq.K8SPodName = fmt.Sprintf("podName-%s", uuid.New().String())

	var requireIPRes proto.IPResMgr2IPAMRequireIPRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&requireIPReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ip/").Post("require").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&requireIPRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ip/require failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		logger.Fatalf("post %sv1/ip/require failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("requireIPRes:%s\n", litter.Sdump(&requireIPRes))
}

func testReleaseIP() {
	var releaseIPReq proto.IPAM2IPResMgrReleaseIPReq
	releaseIPReq.ReqID = ksuid.New().String()
	releaseIPReq.K8SApiResourceKind = proto.K8SApiResourceKindDeployment
	releaseIPReq.K8SClusterID = "cluster-1"
	releaseIPReq.K8SNamespace = "default"
	releaseIPReq.K8SPodName = *unBindPodName

	var releaseIPRes proto.IPResMgr2IPAMReleaseIPRes

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, _ := json.Marshal(&releaseIPReq)

	scli := sling.New().Base(*srvIPResMgrAddr).Set("Content-Type", "text/plain; charset=utf-8")

	res, err := scli.Path("v1/ip/").Post("release").Body(strings.NewReader(calm_utils.Bytes2String(jsonData))).Receive(&releaseIPRes, nil)
	if err != nil {
		logger.Fatalf("post %sv1/ip/release failed. err:%s", *srvIPResMgrAddr, err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		logger.Fatalf("post %sv1/ip/release failed. res.StatusCode:%d", *srvIPResMgrAddr, res.StatusCode)
	}

	logger.Printf("releaseIPRes:%s\n", litter.Sdump(&releaseIPRes))
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	logger = calm_utils.NewSimpleLog(nil)

	switch *testType {
	case 1:
		testCreateIPPool()
	case 2:
		testReleaseIPPool()
	case 3:
		testScaleIPPool()
	case 4:
		wg.Add(*parallel)
		for i := 0; i < *parallel; i++ {
			go func() {
				defer wg.Done()
				testRequireIP()
			}()
		}
	case 5:
		testReleaseIP()
	default:
		logger.Fatalf("Not support type:%d\n", *testType)
	}

	wg.Wait()
	logger.Println("test completed")
}
