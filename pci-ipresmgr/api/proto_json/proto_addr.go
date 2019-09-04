/*
 * @Author: calm.wu
 * @Date: 2019-09-04 18:47:19
 * @Last Modified by:   calm.wu
 * @Last Modified time: 2019-09-04 18:47:19
 */

package protojson

// K8SAddrInfo 分配地址信息
type K8SAddrInfo struct {
	IP                string
	MacAddr           string
	NetRegionalID     string
	SubNetID          string
	PortID            string
	SubNetGatewayAddr string
}
