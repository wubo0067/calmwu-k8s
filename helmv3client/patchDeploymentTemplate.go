/*
 * @Author: calm.wu
 * @Date: 2020-01-08 13:50:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-08 15:24:40
 */

package main

import (
	"bufio"
	"bytes"
	"strings"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

const (
	tagKindStr                                     = "kind: "
	tagDeploymentKindStr                           = "kind: Deployment"
	tagServiceKindStr                              = "kind: Service"
	tagDeploymentSpecStr                           = "spec:"
	tagDeploymentSpecTemplateStr                   = "  template:"
	tagDeploymentSpecTemplateMetadataStr           = "    metadata:"
	tagDeploymentSpecTemplateMetadataAnnotationStr = "      annotations:"

	tagDeploymentMetadataStr                   = "metadata:"
	tagDeploymentMetadataLabelsStr             = "  labels:"
	tagDeploymentSpecSelectorStr               = "  selector:"
	tagDeploymentSpecSelectorMatchlabelsStr    = "    matchLabels:"
	tagDeploymentSpecTemplateMetadataLabelsStr = "      labels:"
)

var (
	sciAnnotations = []string{
		"        io.kubernetes-network.region-id: {{ .Values.Network.RegionID }}",
		"        io.kubernetes.cri.untrusted-workload: \"true\"",
	}

	sciLabels = []string{
		"pci-clusterid: {{ .Values.Label.ClusterID }}",
		"pci-username: {{ .Values.Lable.UserName }}",
	}
)

// 解析deployment template 文件，加上sci相关信息
func patchDeploymentTemplate(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, templateTagKind, error) {
	calm_utils.Debugf("---patchDeploymentTemplate start line:%d---", lineNum)
	tagKind := tagKindNone
	lineContentStr := ""
	bAlreadyRead := false
	bCanRead := scanner.Scan()

	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)
		lineNum++

		// 找到deployment.metadata节点
		if strings.Compare(lineContentStr, tagDeploymentMetadataStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bAlreadyRead = patchLabelsInDeploymentMetaDataRegion(lineNum, scanner, newTemplateBuf)
			calm_utils.Debugf("--->completed deployment.metadata tag node, deployment.metadata end line:%d", lineNum-1)
		}

		// 找到deployment.spec节点，metadata的下一个节点就是
		if strings.Compare(lineContentStr, tagDeploymentSpecStr) == 0 {
			if !bAlreadyRead {
				newTemplateBuf.WriteString(lineContentStr)
				newTemplateBuf.WriteByte('\n')
				lineNum++
			}
			bAlreadyRead = false
			lineNum, lineContentStr, bAlreadyRead, bCanRead = patchInDeploymentSpecRegion(lineNum, scanner, newTemplateBuf)
			if !bCanRead {
				calm_utils.Debugf("--->completed deployment.spec node, deployment.spec end line:%d", lineNum)
			} else {
				calm_utils.Debugf("--->completed deployment.spec node, deployment.spec end line:%d", lineNum-1)
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
			calm_utils.Debug("--->deployment patch completed")
			break
		}
	}

	return lineNum, tagKind, nil
}

func patchInDeploymentSpecRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, bool) {
	calm_utils.Debugf("---deployment.spec start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	bAlreadyRead := false
	//bfindSelector := false  deployment的selector有用户模板自己负责，这个不是后台创建的，需要业务后台来进行关联
	bCanRead := scanner.Scan()

	// 现在开始解析deployment---spec里面的节点
	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到了deployment.spec.selector，神兵的资源这里不添加
		// if strings.Compare(lineContentStr, tagDeploymentSpecSelectorStr) == 0 {
		// 	newTemplateBuf.WriteString(lineContentStr)
		// 	newTemplateBuf.WriteByte('\n')
		// 	lineNum, lineContentStr, bAlreadyRead = patchInDeploymentSpecSelectorRegion(lineNum, scanner, newTemplateBuf)
		// 	calm_utils.Debugf("--->completed deployment.spec.selector node, deployment.spec.selector end line:%d", lineNum-1)
		// 	bfindSelector = true
		// }

		// 找到deployment.spec.template
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateStr) == 0 {
			//if !bAlreadyRead {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			//}
			lineNum, lineContentStr, bAlreadyRead, bCanRead = patchInDeploymentSpecTemplateRegion(lineNum, scanner, newTemplateBuf)
			if !bCanRead {
				calm_utils.Debugf("--->completed deployment.spec.template node, deployment.spec.template end line:-(%d)", lineNum)
			} else {
				calm_utils.Debugf("--->completed deployment.spec.template node, deployment.spec.template end line:+(%d)", lineNum-1)
			}
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)
		// if bRegionEnd || !bCanRead {
		// 	if !bfindSelector {
		// 		// 没有找到deployment.spec.selector
		// 		calm_utils.Debug("--->not find selector in deployment.spec so add")
		// 		newTemplateBuf.WriteString(tagDeploymentSpecSelectorStr)
		// 		newTemplateBuf.WriteByte('\n')
		// 		newTemplateBuf.WriteString(tagDeploymentSpecSelectorMatchlabelsStr)
		// 		newTemplateBuf.WriteByte('\n')
		// 		for _, sciLabel := range sciLabels {
		// 			newTemplateBuf.WriteString("      ")
		// 			newTemplateBuf.WriteString(sciLabel)
		// 			newTemplateBuf.WriteByte('\n')
		// 		}
		// 		bfindSelector = true
		// 	}
		// }

		if isYamlSplitLine(lineContentStr) {
			return lineNum, lineContentStr, false, bCanRead
		}

		if !bAlreadyRead {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
		}
		bAlreadyRead = false

		if bRegionEnd {
			calm_utils.Debug("--->find deployment.spec end")
			break
		}
	}

	return lineNum, lineContentStr, true, bCanRead
}

func patchInDeploymentSpecSelectorRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool) {
	calm_utils.Debugf("---deployment.spec.selector start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	bCanRead := scanner.Scan()

	// 现在开始解析deployment.spec.selector里面的节点
	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.selector.matchLables
		if strings.Compare(lineContentStr, tagDeploymentSpecSelectorMatchlabelsStr) == 0 {
			calm_utils.Debug("--->find matchLables in deployment.spec.selector so patch")
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')

			for _, sciLabel := range sciLabels {
				newTemplateBuf.WriteString("      ")
				newTemplateBuf.WriteString(sciLabel)
				newTemplateBuf.WriteByte('\n')
			}
			continue
		}

		// 不可能没有matchLabels
		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		bRegionEnd = isRegionEnd(lineContentStr, 4)
		if bRegionEnd {
			calm_utils.Debug("--->find deployment.spec.selector.matchLables end")
			break
		}
	}

	return lineNum, lineContentStr, true
}

func patchInDeploymentSpecTemplateRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, bool) {
	calm_utils.Debugf("---deployment.spec.template start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	bAlreadyRead := false
	bCanRead := scanner.Scan()

	// 现在开始解析deployment.spec.template里面的节点
	for ; bCanRead; bCanRead = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.template.metadata
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bAlreadyRead = patchInDeploymentSpecTemplateMetadataRegion(lineNum, scanner, newTemplateBuf)
			calm_utils.Debugf("--->completed deployment.spec.template.metadata node, deployment.spec.template.metadata end line:%d", lineNum-1)
		}

		if isYamlSplitLine(lineContentStr) {
			return lineNum, lineContentStr, false, bCanRead
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)

		// 如果有判断才能这样写
		if !bAlreadyRead {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
		}
		bAlreadyRead = false

		if bRegionEnd {
			calm_utils.Debug("--->find deployment.spec.template end")
			break
		}
	}

	return lineNum, lineContentStr, true, bCanRead
}

func patchInDeploymentSpecTemplateMetadataRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool) {
	calm_utils.Debugf("---deployment.spec.template.metadata start line:%d---", lineNum)
	findAnnotation := false
	findLabels := false
	lineContentStr := ""
	bRegionEnd := false

	for scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.template.metadata
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataAnnotationStr) == 0 {
			calm_utils.Debug("--->find annotation in deployment.spec.template.metadata so patch")
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')

			// 加上sci的annotation
			for _, sciAnno := range sciAnnotations {
				newTemplateBuf.WriteString(sciAnno)
				newTemplateBuf.WriteByte('\n')
			}
			findAnnotation = true
			continue
		}

		// 找到deployment.spec.template.labels
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataLabelsStr) == 0 {
			calm_utils.Debug("--->find annotation in deployment.spec.template.metadata so patch")
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')

			for _, sciLabel := range sciLabels {
				newTemplateBuf.WriteString("        ")
				newTemplateBuf.WriteString(sciLabel)
				newTemplateBuf.WriteByte('\n')
			}
			findLabels = true
			continue
		}

		// 判断是不是deployment.spec.template.metadata的区域结束
		bRegionEnd = isRegionEnd(lineContentStr, 4)
		if bRegionEnd {
			if !findAnnotation {
				// 要加上annotation
				calm_utils.Debug("--->not find annotation in deployment.spec.template.metadata so add")
				newTemplateBuf.WriteString(tagDeploymentSpecTemplateMetadataAnnotationStr)
				newTemplateBuf.WriteByte('\n')
				for _, sciAnno := range sciAnnotations {
					newTemplateBuf.WriteString(sciAnno)
					newTemplateBuf.WriteByte('\n')
				}
				findAnnotation = true
			}

			if !findLabels {
				calm_utils.Debug("--->not find annotation in deployment.spec.template.labels so add")
				newTemplateBuf.WriteString(tagDeploymentSpecTemplateMetadataLabelsStr)
				newTemplateBuf.WriteByte('\n')
				for _, sciLabel := range sciLabels {
					newTemplateBuf.WriteString("        ")
					newTemplateBuf.WriteString(sciLabel)
					newTemplateBuf.WriteByte('\n')
				}
				findLabels = true
			}
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		if bRegionEnd {
			calm_utils.Debug("--->find deployment.spec.template.metadata end")
			break
		}
	}

	return lineNum, lineContentStr, true
}

func patchLabelsInDeploymentMetaDataRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool) {
	calm_utils.Debugf("---deployment.metadata start line:%d---", lineNum)
	findLabels := false
	lineContentStr := ""
	bRegionEnd := false

	for scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到metadata的labels
		if strings.Compare(lineContentStr, tagDeploymentMetadataLabelsStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')

			// 加上sci的metadata---labels
			for _, sciLabel := range sciLabels {
				newTemplateBuf.WriteString("    ")
				newTemplateBuf.WriteString(sciLabel)
				newTemplateBuf.WriteByte('\n')
			}
			findLabels = true
			continue
		}

		// 判断是否到了deployment---metadata区域结束
		bRegionEnd = isRegionEnd(lineContentStr, 0)
		if bRegionEnd {
			calm_utils.Debug("--->find deployment.metadata end")
			// deployment-metadata区域结束都没有labels，需要加上
			if !findLabels {
				// 要加上labels
				newTemplateBuf.WriteString(tagDeploymentMetadataLabelsStr)
				newTemplateBuf.WriteByte('\n')
				for _, sciLabel := range sciLabels {
					newTemplateBuf.WriteString("    ")
					newTemplateBuf.WriteString(sciLabel)
					newTemplateBuf.WriteByte('\n')
				}
			}
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		if bRegionEnd {
			break
		}
	}

	return lineNum, lineContentStr, true
}

// 判断是否是一块区域的结束
func isRegionEnd(lineContent string, spaceCount int) bool {
	// 首先进行替换
	if len(lineContent) > 0 {

		tempLineContent := strings.TrimSpace(lineContent)
		if len(tempLineContent) > 0 && tempLineContent[0] == '{' {
			return false
		}

		for index := 0; index <= spaceCount; index += 2 {
			if lineContent[index] != ' ' {
				return true
			}
		}
	}
	return false
}

// 判断是否是yaml分割符
func isYamlSplitLine(lineContent string) bool {
	return strings.Compare(lineContent, "---") == 0
}

// 判断是不是一个类型的开始
func isKindTag(lineContent string) templateTagKind {
	if strings.Compare(lineContent, tagDeploymentKindStr) == 0 {
		return tagKindDeployment
	}

	if strings.Compare(lineContent, tagServiceKindStr) == 0 {
		return tagKindService
	}

	if strings.HasPrefix(lineContent, tagKindStr) {
		return tagKindOthers
	}
	return tagKindNone
}
