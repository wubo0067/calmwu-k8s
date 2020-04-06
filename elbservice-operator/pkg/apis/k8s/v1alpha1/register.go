// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the k8s v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=k8s.calmwu.org
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "k8s.calmwu.org", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	// 这个在init函数之前，该对象有AddToScheme方法，用来注册自己的schema
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)
