[calmwu@localhost cni-plugin]$ 

go build -v -o bin/calico ./cmd/calico/
go build -v -o bin/calico-ipam ./cmd/calico-ipam/calico-ipam.go
scp bin/calico bin/calico-ipam root@192.168.6.135:/opt/cni/bin