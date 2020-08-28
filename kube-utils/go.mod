module kube-utils

go 1.14

require (
	github.com/containerd/containerd v1.3.6
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.0.0-20191204025024-5ee1b9f4859a // indirect
	golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e // indirect
	google.golang.org/grpc v1.26.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cri-api v0.18.6
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v0.0.0
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
)

replace (
	k8s.io/api => k8s.io/kubernetes/staging/src/k8s.io/api v0.0.0-20200325144952-9e991415386e
	k8s.io/apiextensions-apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver v0.0.0-20200325144952-9e991415386e
	k8s.io/apimachinery => k8s.io/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20200325144952-9e991415386e
	k8s.io/apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiserver v0.0.0-20200325144952-9e991415386e
	k8s.io/cli-runtime => k8s.io/kubernetes/staging/src/k8s.io/cli-runtime v0.0.0-20200325144952-9e991415386e
	k8s.io/client-go => k8s.io/kubernetes/staging/src/k8s.io/client-go v0.0.0-20200325144952-9e991415386e
	k8s.io/cloud-provider => k8s.io/kubernetes/staging/src/k8s.io/cloud-provider v0.0.0-20200325144952-9e991415386e
	k8s.io/cluster-bootstrap => k8s.io/kubernetes/staging/src/k8s.io/cluster-bootstrap v0.0.0-20200325144952-9e991415386e
	k8s.io/code-generator => k8s.io/kubernetes/staging/src/k8s.io/code-generator v0.0.0-20200325144952-9e991415386e
	k8s.io/component-base => k8s.io/kubernetes/staging/src/k8s.io/component-base v0.0.0-20200325144952-9e991415386e
	k8s.io/cri-api => k8s.io/kubernetes/staging/src/k8s.io/cri-api v0.0.0-20200325144952-9e991415386e
	k8s.io/csi-translation-lib => k8s.io/kubernetes/staging/src/k8s.io/csi-translation-lib v0.0.0-20200325144952-9e991415386e
	k8s.io/kube-aggregator => k8s.io/kubernetes/staging/src/k8s.io/kube-aggregator v0.0.0-20200325144952-9e991415386e
	k8s.io/kube-controller-manager => k8s.io/kubernetes/staging/src/k8s.io/kube-controller-manager v0.0.0-20200325144952-9e991415386e
	k8s.io/kube-proxy => k8s.io/kubernetes/staging/src/k8s.io/kube-proxy v0.0.0-20200325144952-9e991415386e
	k8s.io/kube-scheduler => k8s.io/kubernetes/staging/src/k8s.io/kube-scheduler v0.0.0-20200325144952-9e991415386e
	k8s.io/kubectl => k8s.io/kubernetes/staging/src/k8s.io/kubectl v0.0.0-20200325144952-9e991415386e
	k8s.io/kubelet => k8s.io/kubernetes/staging/src/k8s.io/kubelet v0.0.0-20200325144952-9e991415386e
	k8s.io/kubernetes => k8s.io/kubernetes v1.18.0
	k8s.io/legacy-cloud-providers => k8s.io/kubernetes/staging/src/k8s.io/legacy-cloud-providers v0.0.0-20200325144952-9e991415386e
	k8s.io/metrics => k8s.io/kubernetes/staging/src/k8s.io/metrics v0.0.0-20200325144952-9e991415386e
	k8s.io/sample-apiserver => k8s.io/kubernetes/staging/src/k8s.io/sample-apiserver v0.0.0-20200325144952-9e991415386e
)
