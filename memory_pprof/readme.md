### go tool pprof  http://192.168.6.132:30001/debug/pprof/heap 

go tool pprof ./heap

### curl http://192.168.6.132:30001/debug/pprof/trace?seconds=20 > trace

go tool trace ./trace