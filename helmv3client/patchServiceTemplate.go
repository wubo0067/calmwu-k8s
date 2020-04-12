/*
 * @Author: calm.wu
 * @Date: 2020-01-08 14:08:28
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-04-09 10:59:45
 */

// 加上ClusterIP: None，如果有就改为None
// 有type: 直接设置为ClusterIP
// spec如果没有ClusterIP、type，就直接加上

package main

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/pkg/errors"
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

	var err error
	tagKind := tagKindNone
	lineContentStr := ""

	bCanRead := scanner.Scan()

	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)
		lineNum++

		// 找到service.spec节点
		if strings.Compare(lineContentStr, tagServiceSpecStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bCanRead, err = patchInServiceSpecRegion(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, tagKind, errors.Wrap(err, "patchInServiceSpecRegion failed.")
			}
		}

		if !bCanRead {
			break
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		tagKind = isKindTag(lineContentStr)
		if tagKind != tagKindNone {
			// 解析结束
			calm_utils.Debug("--->service patch completed")
			break
		}
	}

	return lineNum, tagKind, nil
}

// clusterIP是默认形式，可以不填写
func patchInServiceSpecRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("---service.spec start line:%d---", lineNum)

	lineContentStr := ""
	bRegionEnd := false
	bFindSpecClusterIPNode := false

	bCanRead := scanner.Scan()

	// 开始解析service.spec节点内容
	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		if strings.HasPrefix(lineContentStr, tagServiceSpecClusterIPPrefixStr) {
			// 替换
			calm_utils.Debug("--->find service.spec.clusterIP, clear this")
			continue
		}

		// 这里可能没有，因为默认clusterIP不需要填写
		if strings.HasPrefix(lineContentStr, tagServiceSpecTypePrefixStr) {
			// 替换
			serviceTypeStr := strings.TrimSpace(lineContentStr)
			calm_utils.Debugf("--->this service.spec.type:[%s] line:%d", serviceTypeStr, lineNum)
			// 判断是不是clusterIP模式
			r, err := regexp.Compile(`type:\s+ClusterIP`)
			if err != nil {
				err = errors.Wrap(err, "regexp Compile type:\\s+ClusterIP failed.")
				calm_utils.Error(err.Error())
				return lineNum, "", true, err
			}

			bMatched := r.MatchString(serviceTypeStr)
			if bMatched {
				calm_utils.Debugf("--->this service.spec.type=ClusterIP, line:%d", lineNum)
				// 设置，为headless
				newTemplateBuf.WriteString(lineContentStr)
				newTemplateBuf.WriteByte('\n')
				newTemplateBuf.WriteString(headlessServiceSpecClusterIPStr)
				newTemplateBuf.WriteByte('\n')
				bFindSpecClusterIPNode = true
			}
			continue
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)
		if bRegionEnd {
			if !bFindSpecClusterIPNode {
				// 加上
				newTemplateBuf.WriteString(headlessServiceSpecClusterIPStr)
				newTemplateBuf.WriteByte('\n')
				newTemplateBuf.WriteString(headlessServiceSpecTypeStr)
				newTemplateBuf.WriteByte('\n')
			}
			calm_utils.Debugf("--->find service.spec end, line:%d", lineNum)
			return lineNum, lineContentStr, true, nil
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}

	if !bFindSpecClusterIPNode {
		// 加上
		newTemplateBuf.WriteString(headlessServiceSpecClusterIPStr)
		newTemplateBuf.WriteByte('\n')
		newTemplateBuf.WriteString(headlessServiceSpecTypeStr)
		newTemplateBuf.WriteByte('\n')
	}

	calm_utils.Debugf("--->find service.spec end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false, nil
}
