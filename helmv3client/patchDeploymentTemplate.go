/*
 * @Author: calm.wu
 * @Date: 2020-01-08 13:50:47
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-03-17 19:20:46
 */

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

const (
	tagKindStr                                                  = "kind: "
	tagDeploymentKindStr                                        = "kind: Deployment"
	tagServiceKindStr                                           = "kind: Service"
	tagDeploymentSpecStr                                        = "spec:"
	tagDeploymentSpecTemplateStr                                = "  template:"
	tagDeploymentSpecTemplateMetadataStr                        = "    metadata:"
	tagDeploymentSpecTemplateMetadataAnnotationStr              = "      annotations:"
	tagDeploymentSpecTemplateSpecStr                            = "    spec:"
	tagDeploymentSpecTemplateSpecContainersStr                  = "      containers:"
	tagDeploymentSpecTemplateSpecContainersResourcesStr         = "        resources:"
	tagDeploymentSpecTemplateSpecContainersResourcesRequestsStr = "          requests:"
	tagDeploymentSpecTemplateSpecContainersResourcesLimitsStr   = "          limits:"
	tagDeploymentSpecTemplateSpecContainersResourcesMemoryStr   = "            memory:"
	tagDeploymentSpecTemplateSpecContainersResourcesCPUStr      = "            cpu:"

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
	var err error
	bCanScan := scanner.Scan()

	for ; bCanScan; bCanScan = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)
		lineNum++

		// 找到deployment.metadata节点
		if strings.Compare(lineContentStr, tagDeploymentMetadataStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			// 因为这里绝不可能是文件的结尾
			lineNum, lineContentStr, bCanScan = patchLabelsInDeploymentMetaDataRegion(lineNum, scanner, newTemplateBuf)
		}

		// 找到deployment.spec节点，metadata的下一个节点就是
		if strings.Compare(lineContentStr, tagDeploymentSpecStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bCanScan, err = patchInDeploymentSpecRegion(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, tagKind, err
			}
		}

		if !bCanScan {
			// 到了文件结尾，跳出
			break
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		tagKind = isKindTag(lineContentStr)
		if tagKind != tagKindNone {
			// 解析结束
			calm_utils.Debug("--->deployment patch completed")
			break
		}
	}

	return lineNum, tagKind, nil
}

func patchInDeploymentSpecRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("---deployment.spec start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	var err error
	//bfindSelector := false  deployment的selector有用户模板自己负责，这个不是后台创建的，需要业务后台来进行关联
	bCanScan := scanner.Scan()

	// 现在开始解析deployment---spec里面的节点
	for ; bCanScan; bCanScan = scanner.Scan() {
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
			lineNum, lineContentStr, bCanScan, err = patchInDeploymentSpecTemplateRegion(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, lineContentStr, bCanScan, err
			}
		}

		if !bCanScan {
			break
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)
		// if bRegionEnd || !bCanScan {
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
			return lineNum, lineContentStr, true, nil
		}

		if bRegionEnd {
			calm_utils.Debugf("--->find deployment.spec end. line:%d", lineNum-1)
			return lineNum, lineContentStr, true, nil
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}

	calm_utils.Debugf("--->find deployment.spec end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false, nil
}

func patchInDeploymentSpecSelectorRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool) {
	calm_utils.Debugf("---deployment.spec.selector start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	bCanScan := scanner.Scan()

	// 现在开始解析deployment.spec.selector里面的节点
	for ; bCanScan; bCanScan = scanner.Scan() {
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

func patchInDeploymentSpecTemplateRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("---deployment.spec.template start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	var err error
	bCanScan := scanner.Scan()

	// 现在开始解析deployment.spec.template里面的节点
	for ; bCanScan; bCanScan = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.template.metadata
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bCanScan, err = patchInDeploymentSpecTemplateMetadataRegion(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, lineContentStr, bCanScan, err
			}
		}

		if !bCanScan {
			break
		}

		// 找到deployment.spec.template.spec
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateSpecStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			lineNum, lineContentStr, bCanScan, err = patchInDeploymentSpecTemplateSpecRegion(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, lineContentStr, bCanScan, err
			}
		}

		if !bCanScan {
			break
		}

		if isYamlSplitLine(lineContentStr) {
			return lineNum, lineContentStr, true, nil
		}

		bRegionEnd = isRegionEnd(lineContentStr, 0)
		if bRegionEnd {
			calm_utils.Debugf("--->find deployment.spec.template end. line:%d", lineNum-1)
			return lineNum, lineContentStr, true, nil
		}

		// 自己范围的都在这里写
		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}

	calm_utils.Debugf("--->find deployment.spec.template end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false, nil
}

func patchInDeploymentSpecTemplateMetadataRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("---deployment.spec.template.metadata start line:%d---", lineNum)
	findAnnotation := false
	findLabels := false
	lineContentStr := ""
	bRegionEnd := false

	bCanScan := scanner.Scan()

	for ; bCanScan; bCanScan = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.template.metadata.annotations
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

		// 找到deployment.spec.template.metadata.labels
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

			calm_utils.Debugf("--->find deployment.spec.template.metadata end. line:%d", lineNum-1)
			return lineNum, lineContentStr, true, nil
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

	}
	calm_utils.Debugf("--->find deployment.spec.template.metadata end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false, nil
}

//
func patchInDeploymentSpecTemplateSpecRegion(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("--->find deployment.spec.template.spec start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	findResources := false
	var err error

	bCanScan := scanner.Scan()

	for ; bCanScan; bCanScan = scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		// 找到deployment.spec.template.spec.containers.resources
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesStr) == 0 {
			newTemplateBuf.WriteString(lineContentStr)
			newTemplateBuf.WriteByte('\n')
			calm_utils.Debugf("--->find resources at line:%d in deployment.spec.template.spec so patch", lineNum)
			findResources = true
			lineNum, lineContentStr, bCanScan, err = patchContainerResources(lineNum, scanner, newTemplateBuf)
			if err != nil {
				return lineNum, lineContentStr, bCanScan, err
			}
		}

		if !bCanScan {
			break
		}

		bRegionEnd = isRegionEnd(lineContentStr, 4)
		if bRegionEnd {
			calm_utils.Debugf("--->find deployment.spec.template.spec end. line:%d",
				lineNum-1)
			// 直接返回
			return lineNum, lineContentStr, true, nil
		}
		// 自己范围内的都是自己写
		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}
	if !findResources {
		// 没有resource配置，报错
		return lineNum, lineContentStr, true, errors.New("deployment.spec.template.spec.containers.resources{request & limits} must be config")
	}
	// 结束行在内部标识
	calm_utils.Debugf("--->find deployment.spec.template.spec end. end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false, nil
}

//
func patchContainerResources(lineNum int, scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, string, bool, error) {
	calm_utils.Debugf("---deployment.spec.template.spec.containers.resources start line:%d---", lineNum)
	lineContentStr := ""
	bRegionEnd := false
	isRequestResource := false
	isLimitResource := false
	var resourceRequestMemory string
	isPatchRequest := false
	requestCpu := 100

	for scanner.Scan() {
		lineContentStr = scanner.Text()
		calm_utils.Debug(lineContentStr)

		lineNum++

		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesRequestsStr) == 0 {
			isRequestResource = true
			isLimitResource = false
			// 这里不写入
			continue
		}

		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesLimitsStr) == 0 {
			isLimitResource = true
			isRequestResource = false
		}

		if strings.Contains(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesCPUStr) && isLimitResource {
			// 这里要判断limit是否有效，而且要计算出request的值
			requestCPUStr := strings.TrimSpace(lineContentStr)
			calm_utils.Debugf("--->this is resource.limits.cpu:[%s], line:%d", requestCPUStr, lineNum)
			cpuRegexp := regexp.MustCompile(`cpu:\s+"([\d]{3,5})m"`)
			params := cpuRegexp.FindStringSubmatch(requestCPUStr)
			//calm_utils.Debugf("match:%v", match)
			for _, param := range params {
				calm_utils.Debugf("param:%s", param)
			}
			limitCPU, _ := strconv.ParseInt(params[1], 10, 32)
			requestCpu = int(limitCPU) / 5
		}

		if strings.Contains(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesCPUStr) && isRequestResource {
			// 这里不写入
			continue
		}

		if strings.Contains(lineContentStr, tagDeploymentSpecTemplateSpecContainersResourcesMemoryStr) && isRequestResource {
			// 这里不写入
			resourceRequestMemory = lineContentStr
			continue
		}

		bRegionEnd = isRegionEnd(lineContentStr, 8)
		if bRegionEnd {
			// 要保证写到结束块的前面
			isPatchRequest = true
			calm_utils.Debugf("--->find deployment.spec.template.spec.containers.resources end. line:%d", lineNum-1)
			newTemplateBuf.WriteString(tagDeploymentSpecTemplateSpecContainersResourcesRequestsStr)
			newTemplateBuf.WriteByte('\n')
			newTemplateBuf.WriteString(resourceRequestMemory)
			newTemplateBuf.WriteByte('\n')
			newTemplateBuf.WriteString(fmt.Sprintf("            cpu: \"%dm\"", requestCpu))
			newTemplateBuf.WriteByte('\n')
			return lineNum, lineContentStr, true, nil
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}

	if !isPatchRequest {
		newTemplateBuf.WriteString(tagDeploymentSpecTemplateSpecContainersResourcesRequestsStr)
		newTemplateBuf.WriteByte('\n')
		newTemplateBuf.WriteString(resourceRequestMemory)
		newTemplateBuf.WriteByte('\n')
		newTemplateBuf.WriteString(fmt.Sprintf("            cpu: \"%dm\"", requestCpu))
		newTemplateBuf.WriteByte('\n')
	}

	// 在这里返回，说明遇到文件尾了
	calm_utils.Debugf("--->find deployment.spec.template.spec.containers.resources end, This is also the end of the file! line:%d", lineNum)

	// 这里行内容返回空字符串
	return lineNum, "", false, nil
}

//
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
			calm_utils.Debugf("--->find deployment.metadata end line:%d", lineNum-1)
			return lineNum, lineContentStr, true
		}

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')
	}
	calm_utils.Debugf("--->find deployment.metadata end, This is also the end of the file! line:%d", lineNum)
	return lineNum, lineContentStr, false
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
