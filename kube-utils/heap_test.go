/*
 * @Author: calmwu
 * @Date: 2020-07-12 00:07:38
 * @Last Modified by: calmwu
 * @Last Modified time: 2020-07-12 00:22:47
 */

package main

import (
	"log"
	"testing"

	"k8s.io/client-go/tools/cache"
)

func testHeapObjectKeyFunc(obj interface{}) (string, error) {
	return obj.(testHeapObject).name, nil
}

type testHeapObject struct {
	name string
	val  interface{}
}

func mkHeapObj(name string, val interface{}) testHeapObject {
	return testHeapObject{name: name, val: val}
}

func compareInts(val1 interface{}, val2 interface{}) bool {
	first := val1.(testHeapObject).val.(int)
	second := val2.(testHeapObject).val.(int)
	return first < second
}

func TestHeapList(t *testing.T) {
	h := cache.NewHeap(testHeapObjectKeyFunc, compareInts)
	list := h.List()
	if len(list) != 0 {
		log.Fatalf("expected an empty list")
	}

	items := map[string]int{
		"foo": 10,
		"bar": 1,
		"bal": 30,
		"baz": 11,
		"faz": 30,
	}
	for k, v := range items {
		h.Add(mkHeapObj(k, v))
	}
	list = h.List()
	if len(list) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(list))
	}
	for index, obj := range list {
		heapObj := obj.(testHeapObject)
		log.Printf("index:%d key:%s\n", index, heapObj.name)
	}
}
