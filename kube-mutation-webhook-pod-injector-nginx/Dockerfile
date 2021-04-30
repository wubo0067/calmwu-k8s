FROM alpine:latest AS final

# 变量重新声明，必须如此
#ARG app

#COPY config/config.json /dynamic-informer-watchres/config/config.json
COPY bin/nginx-injector-pod-webhook-server /nginx-injector/bin/nginx-injector-pod-webhook-server
RUN chmod +x /nginx-injector/bin/nginx-injector-pod-webhook-server

WORKDIR /nginx-injector/bin

ENTRYPOINT ["./nginx-injector-pod-webhook-server"]