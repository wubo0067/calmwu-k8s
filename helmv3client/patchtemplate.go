/*
 * @Author: calm.wu
 * @Date: 2020-01-04 15:17:45
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-01-08 14:11:33
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"

	calm_utils "github.com/wubo0067/calmwu-go/utils"
)

type templateTagKind int

const (
	tagKindNone templateTagKind = iota
	tagKindOthers
	tagKindDeployment
	tagKindService
)

// type patchKind int

// const (
// 	patchKindDeployment patchKind = iota + 1
// 	patchKindService
// 	patchKindOthers
// )

func patchTemplateFileWithSCISpecific(fileName string) {
	calm_utils.Debugf("\n----------------loadTemplateFile----------------")
	templateData, err := ioutil.ReadFile(fileName)
	if err != nil {
		calm_utils.Fatalf("read file:%s failed. err:%s", fileName, err.Error())
	}

	currTagKind := tagKindNone
	newTemplateBuf := &bytes.Buffer{}

	lineNum := 0
	lineContentStr := ""
	scanner := bufio.NewScanner(bytes.NewBuffer(templateData))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// 读取一行
		lineContentStr = scanner.Text()
		lineNum++

		newTemplateBuf.WriteString(lineContentStr)
		newTemplateBuf.WriteByte('\n')

		currTagKind = isKindTag(lineContentStr)

	KINDCHECK:
		switch currTagKind {
		case tagKindDeployment:
			lineNum, currTagKind, _ = patchDeploymentTemplate(lineNum, scanner, newTemplateBuf)
			calm_utils.Debugf("---->currTagKind:%d<----", currTagKind)
			if currTagKind != tagKindOthers {
				goto KINDCHECK
			}
		case tagKindService:
			lineNum, currTagKind, _ = patchServiceTemplate(lineNum, scanner, newTemplateBuf)
			calm_utils.Debugf("---->currTagKind:%d<----", currTagKind)
			if currTagKind != tagKindOthers {
				goto KINDCHECK
			}
		default:
		}
	}

	ioutil.WriteFile(fmt.Sprintf("%s.patch", fileName), newTemplateBuf.Bytes(), 0755)
	//calm_utils.Debugf("after patch:\n%s", newTemplateBuf.String())
}
