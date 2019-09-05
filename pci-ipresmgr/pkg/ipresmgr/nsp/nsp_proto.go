/*
 * @Author: calm.wu
 * @Date: 2019-09-05 13:52:03
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-05 15:12:41
 */

package nsp

// NSPFixedIP 子网
type NSPFixedIP struct {
	SubnetID string `json:"subnet_id,omitempty" mapstructure:"subnet_id"`   // 子网IDb
	IP       string `json:"ip_address,omitempty" mapstructure:"ip_address"` // IP地址
}

// NSPAllocPort 请求port的结构
type NSPAllocPort struct {
	NetRegionalID string        `json:"network_id" mapstructure:"network_id"`         // 网络域ID
	DeviceID      string        `json:"device_id" mapstructure:"device_id"`           // 填写调用方的自定义信息，方便定位
	DeviceOwner   string        `json:"device_owner" mapstructure:"device_owner"`     // 填写compute:kata
	Name          string        `json:"name" mapstructure:"name"`                     // 填写调用方的自定义信息，方便定位
	AdminStateUp  bool          `json:"admin_state_up" mapstructure:"admin_state_up"` // true
	FixedIPs      []*NSPFixedIP `json:"fixed_ips" mapstructure:"fixed_ips"`
}

// NSPAllocPortsReq 分配ports
type NSPAllocPortsReq struct {
	PortLst []*NSPAllocPort `json:"ports" mapstructure:"ports"`
}

// NSPAllocPortResult 分配的port结果
type NSPAllocPortResult struct {
	Name       string       `json:"name" mapstructure:"name"`               // 填写调用方的自定义信息，方便定位
	MacAddress string       `json:"mac_address" mapstructure:"mac_address"` //
	FixedIPs   []NSPFixedIP `json:"fixed_ips" mapstructure:"fixed_ips"`     // 返回IP和子网
	PortID     string       `json:"id" mapstructure:"id"`                   // 这个是PortID
}

// NSPAllocPortsRes 分配的结果
type NSPAllocPortsRes struct {
	PortLst []NSPAllocPortResult `json:"ports" mapstructure:"ports"` // 结果
}
