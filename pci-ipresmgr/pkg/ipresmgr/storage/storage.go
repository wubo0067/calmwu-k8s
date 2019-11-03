/*
 * @Author: calm.wu
 * @Date: 2019-08-29 18:44:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-07 16:30:35
 */

package storage

import (
	proto "pci-ipresmgr/api/proto_json"
)

// StoreMgr 存储接口
type StoreMgr interface {
	// Start 启动存储管理
	Start(Option) error

	// Stop 停止存储管理
	Stop()

	// Register 注册自己，保证实例id唯一
	Register(listenAddr string, listenPort int) error

	// UnRegister 注销自己
	UnRegister()

	// CheckRecycledResources 检查对应资源是否存在，bool = true存在，int=副本数量，
	CheckRecycledResources(k8sResourceID string) (bool, int, error)

	// SetAddrInfosToK8SResourceID 为k8s资源设置地址资源
	SetAddrInfosToK8SResourceID(K8SResourceID string, k8sResourceType proto.K8SApiResourceKindType, k8sAddrInfos []*proto.K8SAddrInfo) error

	// BindAddrInfoWithK8SPodUniqueName 获取一个地址信息，和k8s资源绑定
	BindAddrInfoWithK8SPodUniqueName(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, podUniqueName string) (*proto.K8SAddrInfo, error)

	// UnbindAddrInfoWithK8SPodID 地址和k8s资源解绑
	UnbindAddrInfoWithK8SPodID(k8sResourceType proto.K8SApiResourceKindType, podUniqueName string) error

	// AddK8SResourceAddressToRecycle 加入回收站，待租期到期回收
	AddK8SResourceAddressToRecycle(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType) error

	// SetJobNetInfo 设置job、cronjob的网络信息
	SetJobNetInfo(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, netRegionalID, subNetID,
		subNetGatewayAddr string, subNetCIDR string) error

	// GetJobNetInfo 查询Job、Cronjob的网络信息, 网络域id， 子网id，子网网关地址, subnetcidr
	GetJobNetInfo(k8sResourceID string) (string, string, string, string, error)

	// DelJobNetInfo 删除Job、Cronjob的网络信息
	DelJobNetInfo(k8sResourceID string) error

	// BindJobPodWithPortID 绑定job、cronjob的podid和网络地址
	BindJobPodWithPortID(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, podIP string, portID string, podUniqueName string) error

	// UnbindJobPodWithPortID 解绑job、cronjob的podid和网络地址
	UnbindJobPodWithPortID(podUniqueName string) error

	// ReduceK8SResourceAddrs 给k8s资源地址数量进行缩容
	ReduceK8SResourceAddrs(k8sResourceID string, reduceCount int) error

	// AddScaleDownMarked 添加缩减标记，每个一条记录
	AddScaleDownMarked(k8sResourceID string, k8sResourceType proto.K8SApiResourceKindType, originalReplicas int, scaleDownSize int) error

	// QueryK8SResourceKindByPodUniqueName 在表tbl_K8SResourceIPBind查询pod对应的k8s类型，查不到就是job和cronjob
	QueryK8SResourceKindByPodUniqueName(podUniqueName string) proto.K8SApiResourceKindType
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
