# Shortcut targets
default: build

## Build binary for current platform
all: build

BUILDARCH ?= $(shell uname -m)

# canonicalized names for host architecture
ifeq ($(BUILDARCH),aarch64)
        BUILDARCH=arm64
endif
ifeq ($(BUILDARCH),x86_64)
        BUILDARCH=amd64
endif

# unless otherwise set, I am building for my own architecture, i.e. not cross-compiling
ARCH ?= $(BUILDARCH)

# canonicalized names for target architecture
ifeq ($(ARCH),aarch64)
        override ARCH=arm64
endif
ifeq ($(ARCH),x86_64)
        override ARCH=amd64
endif

RACECHECK = -race
ifeq ($(env), product)
RACECHECK =
endif

SRCFILES=$(shell find pkg cmd api table -name '*.go')

#编译的目录
BIN=bin/$(ARCH)
MAKE_SURE_BIN_EXIST := $(shell mkdir -p $(BIN))

.PHONY: clean
clean:
	rm -rf $(BIN)
#	go get -u ./...

###############################################################################
# Building the binary
###############################################################################
BUILDTIME   =   `date +%FT%T%z`
VERSION     ?= $(shell git describe --always --dirty)
LDFLAGS     =   -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME) -linkmode external -extldflags -static"
GCFLAGS     =   -gcflags "-N -l -m"
GOTAGS      :=  -tags static -tags deadlock

build: $(BIN)/ipresmgr-srv

$(BIN)/ipresmgr-srv : $(SRCFILES)
	gofmt -l -w $(SRCFILES)
	CGO_ENABLED=1 go build $(RACECHECK) -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $@ ./cmd/ipresmgr/ipresmgr.go
