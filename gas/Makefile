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
LDFLAGS     =   -ldflags "-X main.version=${VERSION} -X main.buildtime=${BUILDTIME}" 
LDFLAGS    +=   -ldflags '-linkmode "external" -extldflags "-static"'
GCFLAGS     =   -gcflags "-N -l -m"
GOTAGS      :=  -tags static -tags deadlock

build: proto $(BIN)/gas

proto: 
	make -C api/protobuf

web-hello:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/web/web-hello ./cmd/web/hello	

srv-usermgr:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-usermgr ./cmd/srv/usermgr	

srv-stringprocess:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-stringprocess ./cmd/srv/stringprocess	

pci-gateway:
	CGO_ENABLED=1 GOOS=linux goarch=amd64 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/gateway/pci-gateway ./cmd/gateway

pci-srv-k8soperator:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/k8s-operator ./cmd/srv/k8s-operator

$(BIN)/gas : $(SRCFILES)
	gofmt -l -w $(SRCFILES)
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/web/web-hello ./cmd/web/hello
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-stringprocess ./cmd/srv/stringprocess
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-usermgr ./cmd/srv/usermgr
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/gateway/pci_gateway ./cmd/gateway
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/k8s-operator ./cmd/srv/k8s-operator