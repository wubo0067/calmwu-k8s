/*
 * @Author: calm.wu
 * @Date: 2019-07-11 14:22:14
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-09-04 14:50:31
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
	// IPResMgrErrnoGetIPFailed 获取IP失败
	IPResMgrErrnoGetIPFailed
)
