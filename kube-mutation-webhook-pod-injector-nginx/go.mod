module pod.injector.nginx

go 1.16

require (
	github.com/emirpasic/gods v1.12.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.7.1
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529
	github.com/json-iterator/go v1.1.10
	github.com/pkg/errors v0.9.1
	github.com/sanity-io/litter v1.5.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	istio.io/istio v0.0.0-20210430230726-22a21763af0f
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
)

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible

replace github.com/docker/docker => github.com/moby/moby v20.10.6+incompatible
