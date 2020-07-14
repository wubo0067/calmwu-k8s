/*
 * @Author: calmwu
 * @Date: 2020-07-11 17:43:42
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-07-12 00:18:25
 */

package main

import (
	"fmt"
	"testing"
)

func TestParseEndpoint(t *testing.T) {
	for index, endPoint := range DefaultRuntimeEndpoints {
		fmt.Printf("index:%d endPoint:%s\n", index, endPoint)

		protocol, addr, err := ParseEndpoint(endPoint)
		if err == nil {
			t.Logf("endPoint:%s scheme:[%s] path[%s]", endPoint, protocol, addr)
		} else {
			t.Errorf("endPoint:%s parse failed, err:%s", endPoint, err.Error())
		}
	}
}
