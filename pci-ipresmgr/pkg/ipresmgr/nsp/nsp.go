/*
 * @Author: calm.wu
 * @Date: 2019-09-04 17:31:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-05 15:28:01
 */

// Package nsp 处理和nsp的交互
package nsp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	proto "pci-ipresmgr/api/proto_json"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

// NSPMgrItf nsp交互接口
type NSPMgrItf interface {
	// AllocAddrResources 从nsp获取资源
	AllocAddrResources(k8sResourceID string, replicas int, netRegionID string, subNetID string, subNetGatewayAddr string, netMask int) ([]*proto.K8SAddrInfo, error)

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
func (ni *nspMgrImpl) AllocAddrResources(k8sResourceID string, replicas int, netRegionID string,
	subNetID string, subNetGatewayAddr string, netMask int) ([]*proto.K8SAddrInfo, error) {
	k8sAddrInfoLst := make([]*proto.K8SAddrInfo, 0)

	// TODO: 失败了要回滚
	var rollBackFlag int
	defer func(flag *int) {
		if *flag == 1 {
			for _, k8sAddrInfo := range k8sAddrInfoLst {
				calm_utils.Infof("Rollback PortID:%s", k8sAddrInfo.PortID)
				ni.ReleaseAddrResources(k8sAddrInfo.PortID)
			}
		}
	}(&rollBackFlag)

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
				NetRegionalID: netRegionID,
				DeviceID:      k8sResourceID,
				DeviceOwner:   "compute:kata",
				Name:          k8sResourceID,
				AdminStateUp:  true,
				FixedIPs: []*NSPFixedIP{
					&NSPFixedIP{
						SubnetID: subNetID,
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
			rollBackFlag = 1
			return nil, err
		}
		// 发送请求
		resData, httpCode, err := calm_utils.PostRequest(ni.nspURL, serialData)
		if err != nil {
			err = errors.Wrapf(err, "Post k8sResourceID:[%s] allocPortsURL:[%s] failed.", k8sResourceID, ni.nspURL)
			calm_utils.Error(err)
			rollBackFlag = 1
			return nil, err
		}

		if httpCode >= 200 && httpCode <= 299 {
			// 解析请求
			allocPortsRes := &NSPAllocPortsRes{}
			err = json.Unmarshal(resData, allocPortsRes)
			if err != nil {
				err = errors.Wrapf(err, "Unmarshal k8sResourceID:[%s] allocPortsRes failed.", k8sResourceID)
				calm_utils.Error(err)
				rollBackFlag = 1
				return nil, err
			}

			for index := range allocPortsRes.PortLst {
				portResult := &allocPortsRes.PortLst[index]

				calm_utils.Debugf("%d portResult:%+v", index, portResult)
				k8sAddrInfoLst = append(k8sAddrInfoLst, &proto.K8SAddrInfo{
					IP:                fmt.Sprintf("%s/%d", portResult.FixedIPs[0].IP, netMask),
					MacAddr:           portResult.MacAddress,
					NetRegionalID:     netRegionID,
					SubNetID:          subNetID,
					PortID:            portResult.PortID,
					SubNetGatewayAddr: subNetGatewayAddr,
				})
			}
			replicas -= batchCount
		} else {
			err = errors.Errorf("Post to %s, httpCode is not 2xx", ni.nspURL)
			calm_utils.Errorf(err.Error())
			return nil, err
		}
	}
	return k8sAddrInfoLst, nil
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
		calm_utils.Debugf("Post to %s, httpCode:%d", delPortURL, res.StatusCode)
		defer res.Body.Close()
		ioutil.ReadAll(res.Body)
	}

	return nil
}

// NSPInit 初始化
func NSPInit(nspUrl string) {
	initOnce.Do(func() {
		NSPMgr = &nspMgrImpl{
			httpClient: calm_utils.NewBaseHttpClient(2*time.Second, 6, 2),
			nspURL:     nspUrl,
		}
	})
}
