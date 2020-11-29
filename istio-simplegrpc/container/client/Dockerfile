FROM alpine:latest AS final
WORKDIR /istio-simplegrpc
COPY istio_simplegrpc_client ./bin/

ENTRYPOINT ["/istio-simplegrpc/bin/istio_simplegrpc_client"]