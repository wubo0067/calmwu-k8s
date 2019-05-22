FROM centos

RUN mkdir -p /etc/kubernetes/pki
ADD bin/k8s_webhook_calm /k8s_webhook_calm
ENTRYPOINT ["/k8s_webhook_calm"]