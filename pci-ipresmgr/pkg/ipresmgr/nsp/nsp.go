/*
 * @Author: calm.wu
 * @Date: 2019-09-04 17:31:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 18:51:42
 */

package nsp

import (
	proto "pci-ipresmgr/api/proto_json"
)

// NSPMgr nsp交互接口
type NSPMgr interface {
	// AllocAddrResources 从nsp获取资源
	AllocAddrResources(netRegionalID string, subNetID string) ([]*proto.K8SAddrInfo, error)

	// ReleaseAddrResources 释放资源
	ReleaseAddrResources(portID string)
}
