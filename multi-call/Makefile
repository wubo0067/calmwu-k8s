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

build: proto $(BIN)/multi-call

proto:
	make -C api/protobuf
	protoc --proto_path=.:/home/calm/Program/go_space/src:../ --go_out=. --micro_out=. api/protobuf/common/sci_error/error.proto
	protoc --proto_path=.:/home/calm/Program/go_space/src:../ --go_out=. --micro_out=. api/protobuf/sci_auth/sci_openapiauth.proto

srv-strsplit:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-strsplit ./cmd/srv/strsplit

srv-strupper:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-strupper ./cmd/srv/strupper

cli-strproc:
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/client/cli-strproc ./cmd/client/strproc

$(BIN)/multi-call : $(SRCFILES)
	gofmt -l -w $(SRCFILES)
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-strsplit ./cmd/srv/strsplit
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/srv/srv-strupper ./cmd/srv/strupper
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $(BIN)/client/cli-strproc ./cmd/client/strproc

