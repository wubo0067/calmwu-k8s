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

// NetELBInstance 网络ELB实例
type NetELBInstance struct {
	VpcID             string `json:"vpcid"`
	NetworkID         string `json:"networkid"`         // 网络域id
	SubnetID          string `json:"subnetid"`          // 子网id
	RegionID          string `json:"regionid"`          // 实例所在地域id
	AutoRenew         bool   `json:"autorenew"`         // 是否自动续费，计费类型为包月
	ZoneID            string `json:"zoneid"`            // 实例可用区id
	DisplayName       string `json:"displayname"`       // 实例的显示名
	LoadbalanceTypeID string `json:"loadbalancetypeid"` // 产品规格
	LoadType          string `json:"loadtype"`          // 负载均衡类型F5 LVX LVS
	VIP               string `json:"vip"`               // vip
	AccessType        string `json:"accesstype"`        // 访问类型，INSIDE OUTSIDE
	ELBInstanceID     string `json:"elbinstanceid"`     // elbinstanceid
}

// ELBListener 监听器的参数
type ELBListener struct {
	DidplayName   string `json:"DidplayName"`
	FrontPort     int32  `json:"FrontPort"`
	Protocol      string `json:"protocol"`
	LBStrategy    string `json:"lbstrategy"`    // WRR加权轮询、WLC加权最少连接
	ContainerPort int32  `json:"containerport"` // 后端容器端口
	ListenerID    string `json:"listenerid"`    // 监听器id
	PoolID        string `json:"poolid"`        // 对应的poolid
	FrontProtocol string `json:"frontprotcol"`  // 前端接入协议
}

// ELBServiceSpec defines the desired state of ELBService
type ELBServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ElbInstance NetELBInstance    `json:"elbinstance"`
	Listeners   []ELBListener     `json:"elblisteners"`
	Selector    map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`
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
	Phase          ELBServicePhase `json:"phase"`
	Reason         string          `json:"reason,omitempty"`
	PodCount       int32           `json:"podcount,omitempty"`
	PodInfos       []ELBPodInfo    `json:"podinfos,omitempty"`
	LastUpdateTime metav1.Time     `json:"lastUpdateTime"`
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
