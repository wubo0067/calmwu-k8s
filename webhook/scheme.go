/*
 * @Author: calm.wu
 * @Date: 2019-05-22 14:54:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-22 14:55:23
 */
package main
import (
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	runtimeScheme         = runtime.NewScheme()
	codecs                = serializer.NewCodecFactory(runtimeScheme)
	universalDeserializer = codecs.UniversalDeserializer()
)

func init() {
	addToScheme(runtimeScheme)
}

func addToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(admissionv1beta1.AddToScheme(scheme))
	utilruntime.Must(admissionregistrationv1beta1.AddToScheme(scheme))
}

