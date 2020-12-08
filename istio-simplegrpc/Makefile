APP_PREFIX := istio_simplegrpc
PORT = 8080
RELEASE ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all clean run kube server client proto container

all: proto client server

proto:
	protoc -I${GOOGLEAPIS_DIR} -I. --include_imports --include_source_info --descriptor_set_out=proto/helloworld/helloworld.Greeter.pd --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/helloworld/helloworld.proto

client:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=vendor \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
		-o container/client/${APP_PREFIX}_client src/client/client.go	

server:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=vendor \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
		-o container/server/${APP_PREFIX}_server src/server/server.go		

container:
	crictl rmi ${APP_PREFIX}_client:${RELEASE}
	crictl rmi ${APP_PREFIX}_server:${RELEASE}
	docker build -t ${APP_PREFIX}_client:${RELEASE} container/client		
	docker build -t ${APP_PREFIX}_server:${RELEASE} container/server
	docker save ${APP_PREFIX}_server:${RELEASE} -o ${APP_PREFIX}_server:${RELEASE}.tar
	docker save ${APP_PREFIX}_client:${RELEASE} -o ${APP_PREFIX}_client:${RELEASE}.tar
	ctr -n=k8s.io images import ${APP_PREFIX}_server:${RELEASE}.tar
	ctr -n=k8s.io images import ${APP_PREFIX}_client:${RELEASE}.tar
	-docker rmi ${APP_PREFIX}_client:${RELEASE}
	-docker rmi ${APP_PREFIX}_server:${RELEASE}
	rm ${APP_PREFIX}_server:${RELEASE}.tar ${APP_PREFIX}_client:${RELEASE}.tar
 
clean:
	find ./proto -regex ".*\.go\|.*\.pd"|xargs rm
	rm -f container/server/${APP_PREFIX}_server
	rm -f container/client/${APP_PREFIX}_client
	-docker rmi ${APP_PREFIX}_client:${RELEASE}
	-docker rmi ${APP_PREFIX}_server:${RELEASE}	