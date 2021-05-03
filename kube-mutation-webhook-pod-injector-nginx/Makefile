APP := nginx-injector-pod-webhook-server
RELEASE ?= v0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
WORKDIR=$(shell pwd)
PACKAGES=`go list ./... | grep -v /vendor/`
GOFILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")

.PHONY: all build container-build image clean list lint

all: build

list:
	@echo ---------list source---------
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx golang:1.16 go list ./... | grep -v /vendor/
	@echo ---------list package---------
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx golang:1.16 find . -name "*.go" -type f -not -path "./vendor/*"

container-build:
	# 删除失败，中间镜像
	-docker images | grep "<none>" | awk '{print $3}' | xargs docker rmi
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.16 \
		go build -mod=vendor -ldflags '-X main.Release=${RELEASE} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}' -gcflags 'all=-N -l' \
		-o bin/${APP} main.go

build:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags '-X main.Release=${RELEASE} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -linkmode external -extldflags -static' -gcflags 'all=-N -l'  -o bin/${APP} cmd/main.go

image:
	-docker rmi littlebull/${APP}:$(RELEASE)
	docker build --rm -t littlebull/${APP}:$(RELEASE) --build-arg commitid=${COMMIT} --build-arg app=${APP} --build-arg release=${RELEASE} --build-arg buildtime=${BUILD_TIME} .
	-docker images | grep "<none>" | awk '{print $3}' | xargs docker rmi

push:
	docker image push littlebull/${APP}:$(RELEASE)

test:
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx golang:1.16 go test -v ./...

fmt:
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx golang:1.16 gofmt -s -w ${GOFILES}

lint:
	docker run --rm -v ${WORKDIR}:/src/kube-mutation-webhook-pod-injector-nginx -w /src/kube-mutation-webhook-pod-injector-nginx golangci/golangci-lint:v1.38.0 \
		golangci-lint run -v --config=.golangci.yml ./...

clean:
	-docker rmi ${APP}:$(RELEASE)
	-rm ${APP}