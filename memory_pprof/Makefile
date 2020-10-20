APP := memory_pprof
RELEASE ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all clean build gcrun

all: build

clean:
	rm -f ${APP}

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=vendor \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
		-x -gcflags 'all=-N -l' -o ${APP} main.go

gcrun:
	GODEBUG=gctrace=1 ./${APP}