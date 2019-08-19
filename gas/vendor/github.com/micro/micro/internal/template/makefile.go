package template

var (
	Makefile = `
GOPATH:=$(shell go env GOPATH)

{{if ne .Type "web"}}
.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/{{.Alias}}/{{.Alias}}.proto

.PHONY: build
build: proto
{{else}}
.PHONY: build
build:
{{end}}
	go build -o {{.Alias}}-{{.Type}} main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t {{.Alias}}-{{.Type}}:latest
`
)
