/*
 * @Author: calm.wu
 * @Date: 2019-07-11 14:22:14
 * @Last Modified by:   calm.wu
 * @Last Modified time: 2019-07-11 14:22:14
 */

package protojson

// 错误码类型
type IPResMgrErrorCode int

const (
	IPResMgrErrnoSuccessed IPResMgrErrorCode = iota
	IPResMgrErrnoCreateIPPoolFailed
	IPResMgrErrnoReleaseIPPoolFailed
	IPResMgrErrnoScaleIPPoolFailed
)
