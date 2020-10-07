APP := calm-k8sapp
PORT = 8080
RELEASE ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all clean build run kube

all: build

clean:
	rm -f ${APP}

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod=vendor \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
		-o ${APP} k8sapp.go 

run: container
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(APP):$(RELEASE)

container: build
	docker build -t ${APP}:${RELEASE} .

kube:
	for t in $(shell find ./deploy -type f -name "*.yaml"); do \
cat $$t | \
    sed -E "s/\{\{(\s*)\.Release(\s*)\}\}/$(RELEASE)/g" | \
    sed -E "s/\{\{(\s*)\.ServiceName(\s*)\}\}/$(APP)/g"; \
echo ---; \
done > tmp.yaml