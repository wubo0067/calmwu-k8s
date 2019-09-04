/*
 * @Author: calm.wu
 * @Date: 2019-08-29 18:44:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 17:27:43
 */

package store

import (
	"context"
	proto "pci-ipresmgr/api/proto_json"
)

// StoreMgr 存储接口
type StoreMgr interface {
	// Start 启动存储管理
	Start(context.Context, Option) error

	// Stop 停止存储管理
	Stop()

	// Register 注册自己，保证实例id唯一
	Register(listenAddr string, listenPort int) error

	// UnRegister 注销自己
	UnRegister()

	// CheckRecycledResources 检查对应资源是否存在，bool = true存在，int=副本数量，
	CheckRecycledResources(k8sReousrceID string) (bool, int, error)

	// GetAddrCountByK8SResourceID 根据资源id名，获取k8s资源对应的地址数量
	GetAddrCountByK8SResourceID(k8sReousrceID string) (int, error)

	// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
	SetAddrInfosToK8SResourceID(K8SResourceID string, k8sResourceType proto.K8SApiResourceKindType, k8sAddrInfos []*K8SAddrInfo) error

	// GetAddrInfoByK8SResourceID 获取一个地址信息
	GetAddrInfoByK8SResourceID(k8sReousrceID string) *K8SAddrInfo
}

type K8SAddrInfo struct {
	IP                string
	MacAddr           string
	NetRegionalID     string
	SubNetID          string
	PortID            string
	SubNetGatewayAddr string
}

// StoreOptions 存储的参数
type StoreOptions struct {
	SrvInstID           string
	Addr                string
	User                string
	Passwd              string
	DBName              string
	IdelConnectCount    int
	MaxOpenConnectCount int
	ConnectMaxLifeTime  string
}

// Option 选项修改
type Option func(*StoreOptions)
