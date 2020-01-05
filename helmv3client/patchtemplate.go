/*
 * @Author: calm.wu
 * @Date: 2020-01-04 15:17:45
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-04 22:55:37
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

type templateTagKind int

const (
	tagKindNone templateTagKind = iota
	tagKindOthers
	tagKindDeployment
	tagKindService
)

type patchKind int

const (
	patchKindDeployment patchKind = iota + 1
	patchKindService
	patchKindOthers
)

const (
	tagKindStr                                     = "kind: "
	tagDeploymentKindStr                           = "kind: Deployment"
	tagServiceKindStr                              = "kind: Service"
	tagDeploymentSpecStr                           = "spec:"
	tagDeploymentSpecTemplateStr                   = "  template:"
	tagDeploymentSpecTemplateMetadataStr           = "    metadata:"
	tagDeploymentSpecTemplateMetadataAnnotationStr = "      annotations:"

	tagDeploymentMetadataLabelsStr             = "  labels"
	tagDeploymentSpecSelectorStr               = "  selector"
	tagDeploymentSpecSelectorMatchlabelsStr    = "    matchLabels"
	tagDeploymentSpecTemplateMetadataLablesStr = "      labels"
)

var (
	sciAnnotations = []string{
		"        io.kubernetes-network.region-id: {{ .Values.Network.RegionID }}",
		"        io.kubernetes.cri.untrusted-workload: \"true\"",
	}

	sciLabes = []string{
		"pci-clusterid: {{ .Values.Label.ClusterID }}",
		"pic-username: {{ .Values.Lable.UserName }}",
		"pic.workload.name: {{ .Values.Lable.WorkloadName }}",
		"pic-workload.type: WORKLOAD_DEPLOYMENT",
	}
)

func patchTemplateFileWithSCISpecific(fileName string) {
	calm_utils.Debugf("\n----------------loadTemplateFile----------------")
	templateData, err := ioutil.ReadFile(fileName)
	if err != nil {
		calm_utils.Fatalf("read file:%s failed. err:%s", fileName, err.Error())
	}

	currTagKind := tagKindNone
	currPatchKind := patchKindOthers
	newTemplateBuf := &bytes.Buffer{}

	lineNum := 0
	incLineCount := 0
	scanner := bufio.NewScanner(bytes.NewBuffer(templateData))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// 读取一行
		lineContent := scanner.Bytes()

		newTemplateBuf.Write(lineContent)
		newTemplateBuf.WriteByte('\n')

		switch currPatchKind {
		case patchKindDeployment:
			incLineCount, currTagKind, _ = patchDeploymentTemplateFile(scanner, newTemplateBuf)

			lineNum += incLineCount
			if currTagKind == tagKindDeployment {
				currPatchKind = patchKindDeployment
			} else if currTagKind == tagKindService {
				currPatchKind = patchKindService
			} else {
				currPatchKind = patchKindOthers
			}
		case patchKindService:
			fallthrough
		case patchKindOthers:
			currTagKind = isKindTag(calm_utils.Bytes2String(lineContent))
			if currTagKind == tagKindDeployment {
				currPatchKind = patchKindDeployment
			} else if currTagKind == tagKindService {
				currPatchKind = patchKindService
			}
			lineNum++
		}
	}

	ioutil.WriteFile(fmt.Sprintf("%s.patch", fileName), newTemplateBuf.Bytes(), 0755)
	//calm_utils.Debugf("after patch:\n%s", newTemplateBuf.String())
}

func patchDeploymentTemplateFile(scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) (int, templateTagKind, error) {
	calm_utils.Debugf("---patchDeploymentTemplateFile---")
	lineNum := 0
	tagKind := tagKindNone
	lineSpecTag := -1
	lineSpecTemplateTag := -1
	lineSpecTemplateMetadataTag := -1

	for scanner.Scan() {
		lineContent := scanner.Bytes()

		// 先找到spec---template---metadata
		lineContentStr := calm_utils.Bytes2String(lineContent)
		calm_utils.Debugf("%s", lineContentStr)

		if strings.Compare(lineContentStr, tagDeploymentSpecStr) == 0 {
			newTemplateBuf.Write(lineContent)
			newTemplateBuf.WriteByte('\n')
			lineSpecTag = lineNum
			calm_utils.Debugf("--->find spec yaml node, lineSpecTag:%d", lineSpecTag)
			lineNum++
			continue
		}

		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateStr) == 0 && lineSpecTag > -1 {
			newTemplateBuf.Write(lineContent)
			newTemplateBuf.WriteByte('\n')
			lineSpecTemplateTag = lineNum
			calm_utils.Debugf("--->find template yaml node, lineSpecTemplateTag:%d", lineSpecTemplateTag)
			lineNum++
			continue
		}

		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataStr) == 0 {
			newTemplateBuf.Write(lineContent)
			newTemplateBuf.WriteByte('\n')
			lineSpecTemplateMetadataTag = lineNum
			// 找到sepc---template---metadata节点，开始插入annotation
			if lineSpecTag > -1 && lineSpecTemplateTag > lineSpecTag && lineSpecTemplateMetadataTag > lineSpecTemplateTag {
				calm_utils.Debugf("--->find metadata yaml node, lineSpecTemplateMetadataTag:%d", lineSpecTemplateMetadataTag)
				incLineCount := patchAnnotationInMetaDataRegion(scanner, newTemplateBuf)
				lineNum += incLineCount
			} else {
				lineNum++
			}
			continue
		}

		lineNum++
		newTemplateBuf.Write(lineContent)
		newTemplateBuf.WriteByte('\n')

		tagKind = isKindTag(calm_utils.Bytes2String(lineContent))
		if tagKind != tagKindNone {
			// 解析结束
			calm_utils.Debug("--->deployment patch completed")
			break
		}
	}

	return lineNum, tagKind, nil
}

func patchAnnotationInMetaDataRegion(scanner *bufio.Scanner, newTemplateBuf *bytes.Buffer) int {
	calm_utils.Debugf("---patchAnnotation---")
	lineNum := 0
	findAnnotation := false

	for scanner.Scan() {
		lineContent := scanner.Bytes()

		// 先找到spec---template---metadata
		lineContentStr := calm_utils.Bytes2String(lineContent)
		calm_utils.Debugf("%s", lineContentStr)

		// 找到metadata的annotation
		if strings.Compare(lineContentStr, tagDeploymentSpecTemplateMetadataAnnotationStr) == 0 {
			newTemplateBuf.Write(lineContent)
			newTemplateBuf.WriteByte('\n')

			// 加上sci的annotation
			for _, sciAnno := range sciAnnotations {
				newTemplateBuf.WriteString(sciAnno)
				newTemplateBuf.WriteByte('\n')
			}
			findAnnotation = true
			continue
		}

		// 判断是不是metadata已经结束
		if len(lineContentStr) >= 5 &&
			strings.HasPrefix(lineContentStr, "    ") &&
			lineContentStr[4] != ' ' {
			calm_utils.Debug("--->find metadata end")
			if !findAnnotation {
				// 要加上annotation
				newTemplateBuf.WriteString(tagDeploymentSpecTemplateMetadataAnnotationStr)
				newTemplateBuf.WriteByte('\n')
				for _, sciAnno := range sciAnnotations {
					newTemplateBuf.WriteString(sciAnno)
					newTemplateBuf.WriteByte('\n')
				}
				findAnnotation = true
			}
		}

		lineNum++
		newTemplateBuf.Write(lineContent)
		newTemplateBuf.WriteByte('\n')
	}

	return lineNum
}

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
