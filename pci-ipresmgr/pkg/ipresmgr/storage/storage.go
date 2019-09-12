/*
 * @Author: calm.wu
 * @Date: 2019-08-29 18:44:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:30:35
 */

package storage

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
	CheckRecycledResources(k8sResourceID string) (bool, int, error)

	// GetAddrCountByK8SResourceID 根据资源id名，获取k8s资源对应的地址数量
	GetAddrCountByK8SResourceID(k8sResourceID string) (int, error)

	// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
	SetAddrInfosToK8SResourceID(K8SResourceID string, k8sResourceType proto.K8SApiResourceKindType, k8sAddrInfos []*proto.K8SAddrInfo) error

	// BindAddrInfoWithK8SPodID 获取一个地址信息，和k8s资源绑定
	BindAddrInfoWithK8SPodID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, bindPodID string) (*proto.K8SAddrInfo, error)

	// UnbindAddrInfoWithK8SPodID 地址和k8s资源解绑
	UnbindAddrInfoWithK8SPodID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, unBindPodID string) error

	// AddK8SResourceAddressToRecycle 加入回收站，待租期到期回收
	AddK8SResourceAddressToRecycle(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType) error

	// SetJobNetInfo 设置job、cronjob的网络信息
	SetJobNetInfo(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, netRegionalID, subNetID, subNetGatewayAddr string) error

	// GetJobNetInfo 查询Job、Cronjob的网络信息, 网络域id， 子网id，子网网关地址
	GetJobNetInfo(k8sResourceID string) (string, string, string, error)

	// DelJobNetInfo 删除Job、Cronjob的网络信息
	DelJobNetInfo(k8sResourceID string) error

	// BindJobPodWithPortID 绑定job、cronjob的podid和网络地址
	BindJobPodWithPortID(k8sResourceID string, podIP string, portID string, podID string) error

	// UnbindJobPodWithPortID 解绑job、cronjob的podid和网络地址
	UnbindJobPodWithPortID(k8sResourceID string, podID string) error

	// ScaledK8SResourceAddrs 给k8s资源地址数量进行缩容
	ScaleDownK8SResourceAddrs(k8sResourceID string, scaleDownSize int) error
}

// StoreOptions 存储的参数
type StoreOptions struct {
	SrvInstID           string
	StoreSvrAddr        string
	User                string
	Passwd              string
	DBName              string
	IdelConnectCount    int
	MaxOpenConnectCount int
	ConnectMaxLifeTime  string
}

// Option 选项修改
type Option func(*StoreOptions)
