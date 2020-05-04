package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LoadBalancingStrategy 负载均衡策略
type LoadBalancingStrategy string

const (
	// LBStrategyRR 轮询
	LBStrategyRR LoadBalancingStrategy = "roundrobin"
	// LBStrategyWRR 加权轮询
	LBStrategyWRR LoadBalancingStrategy = "weightedroundrobin"
	// LBStrategyLC 最少链接
	LBStrategyLC LoadBalancingStrategy = "leastconnectons"
	// LBStrategyDH 目标地址散列调度
	LBStrategyDH LoadBalancingStrategy = "destinationhashing"
	// LBStrategySH 源地址散列调度
	LBStrategySH LoadBalancingStrategy = "sourcehashing"
)

// LoadBalancingStrategies 策略列表
var LoadBalancingStrategies []LoadBalancingStrategy

type ELBServicePhase string

const (
	ELBServiceNone        ELBServicePhase = ""
	ELBServiceCreating    ELBServicePhase = "Creating"
	ELBServiceActive      ELBServicePhase = "Active"
	ELBServiceFailed      ELBServicePhase = "Failed"
	ELBServiceTerminating ELBServicePhase = "Terminating"
	ELBServiceUnknown     ELBServicePhase = "Unknown"
)

// CheckLBStrategy 检查策略合法性
func CheckLBStrategy(strategy LoadBalancingStrategy) bool {
	for _, lbName := range LoadBalancingStrategies {
		if lbName == strategy {
			return true
		}
	}
	return false
}

// PodLabel 选择标签
type PodLabel struct {
	Name  string `json:"labelName"`
	Value string `json:"labelVal"`
}

// ELBListener 监听器的参数
type ELBListener struct {
	Name       string                `json:"name"`
	VIP        string                `json:"vip"`
	Port       int32                 `json:"port"`
	Protocol   string                `json:"protocol"`
	LBStrategy LoadBalancingStrategy `json:"lbstrategy"`
}

// ELBServiceSpec defines the desired state of ELBService
type ELBServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Listener ELBListener       `json:"elblistener"`
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`
}

// ELBPodInfo pod信息
type ELBPodInfo struct {
	Name   string          `json:"name"`
	PodIP  string          `json:"podip"`
	Status corev1.PodPhase `json:"status"`
}

// ELBServiceStatus defines the observed state of ELBService
type ELBServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase    ELBServicePhase `json:"phase"`
	PodCount int32           `json:"podcount"`
	PodInfos []ELBPodInfo    `json:"podinfos"`
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

	LoadBalancingStrategies = []LoadBalancingStrategy{
		LBStrategyRR,
		LBStrategyWRR,
		LBStrategyLC,
		LBStrategyDH,
		LBStrategySH,
	}
}
