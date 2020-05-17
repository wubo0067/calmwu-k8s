package apis

import (
	"calmwu.org/elbservice-operator/pkg/apis/k8s/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	// 注册自己schema的注册方法
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
