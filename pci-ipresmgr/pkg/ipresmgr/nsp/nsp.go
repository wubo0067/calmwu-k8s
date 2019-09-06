/*
 * @Author: calm.wu
 * @Date: 2019-09-04 17:31:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-05 15:28:01
 */

package nsp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	proto "pci-ipresmgr/api/proto_json"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// NSPMgrItf nsp交互接口
type NSPMgrItf interface {
	// AllocAddrResources 从nsp获取资源
	AllocAddrResources(k8sResourceID string, ipPoolCreateReq *proto.WB2IPResMgrCreateIPPoolReq) ([]*proto.K8SAddrInfo, error)

	// ReleaseAddrResources 释放资源
	ReleaseAddrResources(portID string) error
}

type nspMgrImpl struct {
	httpClient *http.Client
	nspURL     string
}

var (
	initOnce sync.Once
	// NSPMgr 全局对象
	NSPMgr NSPMgrItf
	//
	_ NSPMgrItf = &nspMgrImpl{}
)

const (
	nspAllocBatchMaxCount = 20
)

// AllocAddrResources 从nsp获取资源
func (ni *nspMgrImpl) AllocAddrResources(k8sResourceID string, ipPoolCreateReq *proto.WB2IPResMgrCreateIPPoolReq) ([]*proto.K8SAddrInfo, error) {
	k8sAddrLst := make([]*proto.K8SAddrInfo, 0)

	replicas := ipPoolCreateReq.K8SApiResourceReplicas
	for replicas > 0 {
		batchCount := func() int {
			if replicas > nspAllocBatchMaxCount {
				return nspAllocBatchMaxCount
			}
			return replicas
		}()

		allocPortsReq := &NSPAllocPortsReq{}
		for index := 0; index < batchCount; index++ {
			allocPortsReq.PortLst = append(allocPortsReq.PortLst, &NSPAllocPort{
				NetRegionalID: ipPoolCreateReq.NetRegionalID,
				DeviceID:      k8sResourceID,
				DeviceOwner:   "compute:kata",
				Name:          k8sResourceID,
				AdminStateUp:  true,
				FixedIPs: []*NSPFixedIP{
					&NSPFixedIP{
						SubnetID: ipPoolCreateReq.SubnetID,
					},
				},
			})
		}

		// 消息序列化
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		serialData, err := json.Marshal(allocPortsReq)
		if err != nil {
			err = errors.Wrapf(err, "Marshal k8sResourceID:[%s] allocPortsReq failed.", k8sResourceID)
			calm_utils.Error(err)
			return nil, err
		}
		// 发送请求

		allocPortsURL := fmt.Sprintf("%s", ni.nspURL)
		resData, err := calm_utils.PostRequest(allocPortsURL, serialData)
		if err != nil {
			err = errors.Wrapf(err, "Post k8sResourceID:[%s] allocPortsURL:[%s] failed.", k8sResourceID, allocPortsURL)
			calm_utils.Error(err)
			return nil, err
		}

		// 解析请求
		allocPortsRes := &NSPAllocPortsRes{}
		err = json.Unmarshal(resData, allocPortsRes)
		if err != nil {
			err = errors.Wrapf(err, "Unmarshal k8sResourceID:[%s] allocPortsRes failed.", k8sResourceID)
			calm_utils.Error(err)
			return nil, err
		}

		for index := range allocPortsRes.PortLst {
			portResult := &allocPortsRes.PortLst[index]

			calm_utils.Debugf("%d portResult:%+v", index, portResult)
			k8sAddrLst = append(k8sAddrLst, &proto.K8SAddrInfo{
				IP:                portResult.FixedIPs[0].IP,
				MacAddr:           portResult.MacAddress,
				NetRegionalID:     ipPoolCreateReq.NetRegionalID,
				SubNetID:          ipPoolCreateReq.SubnetID,
				PortID:            portResult.PortID,
				SubNetGatewayAddr: ipPoolCreateReq.SubnetGatewayAddr,
			})
		}

		replicas -= batchCount
	}
	return k8sAddrLst, nil
}

// ReleaseAddrResources 释放ip资源
func (ni *nspMgrImpl) ReleaseAddrResources(portID string) error {
	delPortURL := fmt.Sprintf("%s/%s", ni.nspURL, portID)
	calm_utils.Debugf("delPortUrl:%s", delPortURL)

	delReq, _ := http.NewRequest("DELETE", delPortURL, nil)
	res, err := ni.httpClient.Do(delReq)
	if err != nil {
		err = errors.Wrapf(err, "DELETE request:%s failed.", delPortURL)
		calm_utils.Error(err)
		return err
	}
	if res != nil {
		defer res.Body.Close()
	}
	ioutil.ReadAll(res.Body)
	return nil
}

// NSPInit 初始化
func NSPInit(nspUrl string) {
	initOnce.Do(func() {
		NSPMgr = &nspMgrImpl{
			httpClient: calm_utils.NewBaseHttpClient(6, 2),
			nspURL:     nspUrl,
		}
	})
}
