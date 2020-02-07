package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ELBListener 监听器的参数
type ELBListener struct {
	Name     string `json:"name"`
	VIP      string `json:"vip"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

// ELBServiceSpec defines the desired state of ELBService
type ELBServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Listener ELBListener `json:"elblistener"`
}

// ELBServiceStatus defines the observed state of ELBService
type ELBServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ELBService is the Schema for the elbservices API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=elbservices,scope=Namespaced
type ELBService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ELBServiceSpec   `json:"spec,omitempty"`
	Status ELBServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ELBServiceList contains a list of ELBService
type ELBServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ELBService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ELBService{}, &ELBServiceList{})
}
