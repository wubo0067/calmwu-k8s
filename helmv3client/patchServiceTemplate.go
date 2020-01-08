/*
 * @Author: calm.wu
 * @Date: 2020-01-08 14:08:28
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-08 15:24:33
 */

// 加上ClusterIP: None，如果有就改为None
// 有type: 直接设置为ClusterIP
// spec如果没有ClusterIP、type，就直接加上

package main

import (
	"bufio"
	"bytes"
	"strings"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

const (
	tagServiceSpecStr                = "spec:"
	tagServiceSpecClusterIPPrefixStr = "  clusterIP:"
	tagServiceSpecTypePrefixStr      = "  type:"
	headlessServiceSpecClusterIPStr  = "  clusterIP: None"
	headlessServiceSpecTypeStr       = "  type: ClusterIP"
)

func patchServiceTemplate(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, templateTagKind, error) {
	calm_utils.Debugf("---patchServiceTemplate start line:%d---", lineNum)

	tagKind := tagKindNone
	lineContentStr := ""
	bAlreadyRead := false
	bCanRead := scanner.Scan()

	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)
		lineNum++

		// 找到service.spec节点
		if strings.Compare(lineContentStr, tagServiceSpecStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bAlreadyRead, bCanRead = patchInServiceSpecRegion(lineNum, scanner, newTemplateBuf)
			if !bCanRead {
				calm_utils.Debugf("--->completed service.spec node, service.spec end line:%d", lineNum)
			} else {
				calm_utils.Debugf("--->completed service.spec node, service.spec end line:%d", lineNum-1)
			}
		}

		if !bAlreadyRead {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
		}
		bAlreadyRead = false

		tagKind = isKindTag(lineContentStr)
		if tagKind != tagKindNone {
			// 解析结束
			calm_utils.Debug("--->service patch completed")
			break
		}
	}

	return lineNum, tagKind, nil
}

func patchInServiceSpecRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, bool) {
	calm_utils.Debugf("---service.spec start line:%d---", lineNum)

	lineContentStr := ""
	bRegionEnd := false
	//bFindSpecClusterIPNode := false
	//bFindSpecTypeNode := false

	newTemplateBuf.WriteString(headlessServiceSpecClusterIPStr)
	newTemplateBuf.WriteByte('\n')
	newTemplateBuf.WriteString(headlessServiceSpecTypeStr)
	newTemplateBuf.WriteByte('\n')

	bCanRead := scanner.Scan()

	// 开始解析service.spec节点内容
	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		if strings.HasPrefix(lineContentStr, tagServiceSpecClusterIPPrefixStr) {
			// 替换
			calm_utils.Debug("--->find clusterIP in service.spec so replace")
			continue
		}

		if strings.HasPrefix(lineContentStr, tagServiceSpecTypePrefixStr) {
			// 替换
			calm_utils.Debug("--->find type in service.spec so replace")
			continue
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)
		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
		if bRegionEnd {
			calm_utils.Debug("--->find service.spec end")
			break
		}
	}

	// if !bCanRead {
	// 	if !bFindSpecClusterIPNode {
	// 		calm_utils.Debug("--->not find clusterIP in service.spec so add")
	// 		newTemplateBuf.WriteString(headlessServiceSpecClusterIPStr)
	// 		newTemplateBuf.WriteByte('\n')
	// 	}

	// 	if !bFindSpecTypeNode {
	// 		calm_utils.Debug("--->not find type in service.spec so add")
	// 		newTemplateBuf.WriteString(headlessServiceSpecTypeStr)
	// 		newTemplateBuf.WriteByte('\n')
	// 	}
	// }
	return lineNum, lineContentStr, true, bCanRead
}
