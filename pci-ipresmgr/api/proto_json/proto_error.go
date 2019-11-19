/*
 * @Author: calm.wu
 * @Date: 2019-07-11 14:22:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-10-04 16:56:25
 */

package protojson

// IPResMgrErrorCode 错误码类型
type IPResMgrErrorCode int

const (
	// IPResMgrErrnoSuccessed 成功
	IPResMgrErrnoSuccessed IPResMgrErrorCode = iota
	// IPResMgrErrnoCreateIPPoolFailed 创建IPPool失败
	IPResMgrErrnoCreateIPPoolFailed
	// IPResMgrErrnoReleaseIPPoolFailed 释放IPPool失败
	IPResMgrErrnoReleaseIPPoolFailed
	// IPResMgrErrnoScaleIPPoolFailed 扩缩容IPPool失败
	IPResMgrErrnoScaleIPPoolFailed
	// IPResMgrErrnoGetIPFailed IPAM获取IP失败
	IPResMgrErrnoGetIPFailed
	// IPResMgrErrnoReleaseIPFailed IPAM释放IP失败
	IPResMgrErrnoReleaseIPFailed
	// IPResMgrErrnoMaintainForceUnbindIPFailed
	IPResMgrErrnoMaintainForceUnbindIPFailed
	// IPResMgrErrnoMaintainForceReleaseK8SResourceIPPoolFailed
	IPResMgrErrnoMaintainForceReleaseK8SResourceIPPoolFailed
	// IPResMgrErrnoMaintainForceReleasePodIPFailed
	IPResMgrErrnoMaintainForceReleasePodIPFailed
)
