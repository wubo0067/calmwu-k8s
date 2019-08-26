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

SRCFILES=$(shell find pkg cmd api -name '*.go')

#编译的目录
BIN=bin/$(ARCH)
MAKE_SURE_BIN_EXIST := $(shell mkdir -p $(BIN))

.PHONY: clean
clean:
	rm -rf $(BIN)

###############################################################################
# Building the binary
###############################################################################
BUILDTIME   =   `date +%FT%T%z`
#VERSION     :=  `git describe --tags`
VERSION     :=  "0.0.1"
LDFLAGS     =   -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME) -linkmode external -extldflags -static"
#LDFLAGS    +=   -linkmode external -extldflags -static"
GCFLAGS     =   -gcflags "-N -l -m"
GOTAGS      :=  -tags static -tags deadlock

build: $(BIN)/ipresmgr-svr

$(BIN)/ipresmgr-svr : $(SRCFILES)
	gofmt -l -w $(SRCFILES)
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/ipresmgr-svr ./cmd/ipresmgr/ipresmgr.go
