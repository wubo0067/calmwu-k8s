FROM scratch
WORKDIR /k8sapp
ENV PORT 8000

COPY calm-k8sapp .
#USER nonroot:nonroot

ENTRYPOINT ["/k8sapp/calm-k8sapp"]