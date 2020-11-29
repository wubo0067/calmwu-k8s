APP_PREFIX := istio_simplegrpc
PORT = 8080
RELEASE ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all clean build run kube server client proto

all: build

proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/helloworld/helloworld.proto

client:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=vendor \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
		-o bin/${APP_PREFIX}_client src/client/client.go		

clean:
	rm -f bin/${APP_PREFIX}_server
	rm -f bin/${APP_PREFIX}_client