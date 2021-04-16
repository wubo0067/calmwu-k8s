# FROM golang:1.16 AS builder

# ARG app
# ARG release
# ARG buildtime
# ARG commitid

# WORKDIR /build

# COPY *.go .
# COPY go.mod .
# COPY go.sum .
# COPY只拷贝目录中的文件，要想拷贝目录，必须dest指定目录
# ADD vendor ./vendor
#ENV GOPROXY https://goproxy.cn
#RUN go mod vendor

# RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor \
# 		-ldflags "-s -w -X main.Release=${release} \
# 		-X main.Commit=${commitid} -X main.BuildTime=${buildtime}" -gcflags 'all=-N -l' \
# 		-o ${app} dyninf-watchres.go

FROM alpine:latest AS final

# 变量重新声明，必须如此
#ARG app

#COPY config/config.json /dynamic-informer-watchres/config/config.json
COPY bin/dynamic-informer-watchres /dynamic-informer-watchres/bin/dynamic-informer-watchres
RUN chmod +x /dynamic-informer-watchres/bin/dynamic-informer-watchres

WORKDIR /dynamic-informer-watchres/bin

ENTRYPOINT ["./dynamic-informer-watchres"]