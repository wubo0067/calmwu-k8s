/*
 * @Author: CALM.WU
 * @Date: 2020-10-08 11:45:22
 * @Last Modified by: CALM.WU
 * @Last Modified time: 2020-10-08 13:39:00
 */

// Package schemegvkrrestmap ...
package schemegvkrrestmap

import (
	"testing"

	"github.com/sanity-io/litter"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func TestSchemeObject2GVK(t *testing.T) {
	gvks, ret, err := scheme.Scheme.ObjectKinds(&appsv1.Deployment{})
	if err != nil {
		t.Fatalf("scheme ObjectKinds for appsv1.Deployment failed. err:%s", err.Error())
	}

	if !ret {
		t.Log("appsv1.Deployment is unversioned")
	}
	for index := range gvks {
		t.Logf("[%d] gvk:{%s}", index, gvks[index].String())
	}

	gvks, ret, err = scheme.Scheme.ObjectKinds(&appsv1.StatefulSet{})
	if err != nil {
		t.Fatalf("scheme ObjectKinds for appsv1.Deployment failed. err:%s", err.Error())
	}

	if !ret {
		t.Log("appsv1.Deployment is unversioned")
	}
	for index := range gvks {
		t.Logf("[%d] gvk:{%s}", index, gvks[index].String())
	}
}

func findGVR(gvk *schema.GroupVersionKind, cfg *rest.Config) (*meta.RESTMapping, error) {

	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}

func TestGVK2RESTMapping(t *testing.T) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		t.Fatalf("build restCfg failed. err:%s", err.Error())
	}

	gvks, _, err := scheme.Scheme.ObjectKinds(&appsv1.Deployment{})
	if err != nil {
		t.Fatalf("scheme ObjectKinds for appsv1.Deployment failed. err:%s", err.Error())
	}

	metaRestMapping, err := findGVR(&schema.GroupVersionKind{
		Group:   gvks[0].Group,
		Version: gvks[0].Version,
		Kind:    gvks[0].Kind,
	}, restCfg)
	if err != nil {
		t.Fatalf("findGVR failed. err:%s", err.Error())
	}

	t.Logf("metaRestMapping:{%s}", litter.Sdump(metaRestMapping))

	// Resource: schema.GroupVersionResource{
	// 	Group: "apps",
	// 	Version: "v1",
	// 	Resource: "deployments",
	//   },
	//   GroupVersionKind: schema.GroupVersionKind{
	// 	Group: "apps",
	// 	Version: "v1",
	// 	Kind: "Deployment",
	//   },
	//   Scope: &meta.restScope{},
	// }}
}
