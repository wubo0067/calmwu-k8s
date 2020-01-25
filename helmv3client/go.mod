module helmclient

go 1.13

replace (
	// github.com/Azure/go-autorest/autorest has different versions for the Go
	// modules than it does for releases on the repository. Note the correct
	// version when updating.
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.9.0
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

	// Kubernetes imports github.com/miekg/dns at a newer version but it is used
	// by a package Helm does not need. Go modules resolves all packages rather
	// than just those in use (like Glide and dep do). This sets the version
	// to the one oras needs. If oras is updated the version should be updated
	// as well.
	github.com/miekg/dns => github.com/miekg/dns v0.0.0-20181005163659-0d29b283ac0f
	gopkg.in/inf.v0 v0.9.1 => github.com/go-inf/inf v0.9.1
	gopkg.in/square/go-jose.v2 v2.3.0 => github.com/square/go-jose v2.3.0+incompatible

	rsc.io/letsencrypt => github.com/dmcgowan/letsencrypt v0.0.0-20160928181947-1847a81d2087
)

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586 // indirect
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/gin-gonic/gin v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/monnand/dhkx v0.0.0-20180522003156-9e5b033f1ac4 // indirect
	github.com/sanity-io/litter v1.2.0
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/snwfdhmp/errlog v0.0.0-20191219134421-4c9e67f11ebc
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/wubo0067/calmwu-go v0.0.0-20191231073022-8cf6b9680e47
	golang.org/x/mod v0.2.0 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/tools v0.0.0-20200121042740-dbc83e6dc05e // indirect
	golang.org/x/tools/gopls v0.1.8-0.20200121042740-dbc83e6dc05e // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	helm.sh/helm/v3 v3.0.0-beta.5.0.20200119220513-a911600fc2d6
	k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	sigs.k8s.io/yaml v1.1.0
)
