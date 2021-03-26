APP := inplaceupdate-pod
RELEASE ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
WORKDIR=$(shell pwd)

.PHONY: all build container clean

all: build

build:
	docker run --rm -v ${WORKDIR}:/src/inplaceupdate -w /src/inplaceupdate -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.16 \
		go build -mod=vendor -ldflags '-X main.Release=${RELEASE} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}' -gcflags 'all=-N -l' \
		-o ${APP} inplaceupdate.go
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags '-X main.Release=0.0.1 -X main.Commit=6fd956f -X main.BuildTime=2021-03-26_14:37:53 -linkmode external -extldflags -static' -gcflags 'all=-N -l'  -o inplaceupdate-pod inplaceupdate.go

container:
	docker build --rm -t ${APP}:$(RELEASE) --build-arg commitid=${COMMIT} --build-arg app=${APP} --build-arg release=${RELEASE} --build-arg buildtime=${BUILD_TIME} .

clean:
	-rm ${APP}	