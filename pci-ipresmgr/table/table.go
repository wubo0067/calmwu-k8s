/*
 * @Author: calm.wu
 * @Date: 2019-08-29 11:48:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-29 17:44:15
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
	K8SResourceType   string         `db:"k8sresource_type"`
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
	SrvInstanceName        string    `db:"srv_instance_name"`
	Replicas               int       `db:"replicas"`
	K8SResourceID          string    `db:"k8sresource_id"`
	CreateTime             time.Time `db:"create_time"`
	NSPResourceReleaseTime time.Time `db:"nspresource_release_time"`
	NetRegionalID          string    `db:"netregional_id"`
	SubNetID               string    `db:"subnet_id"`
	PortID                 string    `db:"port_id"`
	SubNetGatewayAddr      string    `db:"subnetgatewayaddr"`
	NspResources           []byte    `db:"nsp_resources"`
}

// TblK8SResourceIPRecycleHistoryS 地址回收历史表，TblK8SResourceIPRecycleS删除的记录存放到该表
type TblK8SResourceIPRecycleHistoryS struct {
	ID                     uint      `db:"id"`
	K8SResourceID          string    `db:"k8sresource_id"`
	Replicas               int       `db:"replicas"`
	CreateTime             time.Time `db:"create_time"`
	NSPResourceReleaseTime time.Time `db:"nspresource_release_time"`
	NetRegionalID          string    `db:"netregional_id"`
	SubNetID               string    `db:"subnet_id"`
	PortID                 string    `db:"port_id"`
	NspResources           []byte    `db:"nsp_resources"`
}
