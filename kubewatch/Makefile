.PHONY: default build docker-image test stop clean-images clean

BINARY = kubewatch

VERSION=
BUILD=

PKG            = github.com/bitnami-labs/kubewatch
TRAVIS_COMMIT ?= `git describe --tags`
GOCMD          = go
BUILD_DATE     = `date +%FT%T%z`
GOFLAGS        = -mod=vendor
LDFLAGS       := "-X '$(PKG)/cmd.gitCommit=$(TRAVIS_COMMIT)' \
		          -X '$(PKG)/cmd.buildDate=$(BUILD_DATE)'"
LDFLAGS		   =

LDFLAGS     =   -ldflags "-extldflags=-static"
GCFLAGS     =   -gcflags "-N -l -m"
GOTAGS      :=  -tags static -tags deadlock

default: build test

build:
	go build ${GOFLAGS} ${GCFLAGS} ${LDFLAGS} -o kubewatch main.go

docker-image:
	@docker build -t "${BINARY}" .

test:
	"$(GOCMD)" test -race -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i
