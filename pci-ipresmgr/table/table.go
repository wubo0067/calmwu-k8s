/*
 * @Author: calm.wu
 * @Date: 2019-08-29 11:48:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 16:53:40
 */

package table

import (
	"database/sql"
	"time"
)

// TblIPResMgrSrvRegisgerS 服务启动登记表
type TblIPResMgrSrvRegisgerS struct {
	SrvInstanceName string    `db:"srv_instance_name"`
	SrvAddr         string    `db:"srv_addr"`
	RegisterTime    time.Time `db:"register_time"`
}

// TblK8SResourceIPBindS 地址资源绑定表
type TblK8SResourceIPBindS struct {
	K8SResourceID     string         `db:"k8sresource_id"`
	K8SResourceType   int            `db:"k8sresource_type"` // proto.K8SApiResourceKindType
	IP                string         `db:"ip"`
	MacAddr           string         `db:"mac"`
	NetRegionalID     string         `db:"netregional_id"`
	SubNetID          string         `db:"subnet_id"`
	PortID            string         `db:"port_id"`
	SubNetGatewayAddr string         `db:"subnetgatewayaddr"`
	AllocTime         time.Time      `db:"alloc_time"`
	IsBind            int8           `db:"is_bind"`
	BindPodID         sql.NullString `db:"bind_podid"`
	BindTime          time.Time      `db:"bind_time"`
}

// TblK8SResourceIPRecycleS 地址资源回收表
type TblK8SResourceIPRecycleS struct {
	SrvInstanceName string `db:"srv_instance_name"`
	K8SResourceID   string `db:"k8sresource_id"`
	K8SResourceType int    `db:"k8sresource_type"` // proto.K8SApiResourceKindType
	Replicas        int    `db:"replicas"`
	//UnbindCount            int       `db:"unbind_count"`
	CreateTime             time.Time `db:"create_time"`
	NSPResourceReleaseTime time.Time `db:"nspresource_release_time"`
	RecycleObjectID        string    `db:"recycle_object_id"`
	//NetRegionalID          string    `db:"netregional_id"`
	//SubNetID               string    `db:"subnet_id"`
	//PortID                 string    `db:"port_id"`
	//SubNetGatewayAddr      string    `db:"subnetgatewayaddr"`
	//NspResources           []byte    `db:"nsp_resources"`
}

// TblK8SResourceIPRecycleHistoryS 地址回收历史表，TblK8SResourceIPRecycleS删除的记录存放到该表
type TblK8SResourceIPRecycleHistoryS struct {
	ID                     uint      `db:"id"`
	K8SResourceID          string    `db:"k8sresource_id"`
	K8SResourceType        int       `db:"k8sresource_type"` // proto.K8SApiResourceKindType
	Replicas               int       `db:"replicas"`
	CreateTime             time.Time `db:"create_time"`
	NSPResourceReleaseTime time.Time `db:"nspresource_release_time"`
	NetRegionalID          string    `db:"netregional_id"`
	SubNetID               string    `db:"subnet_id"`
	PortID                 string    `db:"port_id"`
	NspResources           []byte    `db:"nsp_resources"`
}

// TblK8SJobNetInfoS Job和CronJob的网络信息
type TblK8SJobNetInfoS struct {
	K8SResourceID     string `db:"k8sresource_id"`
	K8SResourceType   int    `db:"k8sresource_type"` // proto.K8SApiResourceKindType
	NetRegionalID     string `db:"netregional_id"`
	SubNetID          string `db:"subnet_id"`
	SubNetGatewayAddr string `db:"subnetgatewayaddr"`
}

// TblK8SJobIPBindS job 和 cronjob 的pod的ip地址绑定信息
type TblK8SJobIPBindS struct {
	K8SResourceID string `db:"k8sresource_id"`
	IP            string `db:"ip"`
	BindPodID     string `db:"bind_podid"`
	PortID        string `db:"port_id"`
}

type TblK8SScaleDownMarkS struct {
	RecycleMarkID   string         `db:"recycle_mark_id"`
	K8SResourceID   string         `db:"k8sresource_id"`
	K8SResourceType int            `db:"k8sresource_type"` // proto.K8SApiResourceKindType
	PodID           sql.NullString `db:"pod_id"`
}
