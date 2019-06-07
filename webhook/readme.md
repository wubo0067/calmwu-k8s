docker build -t="littlebull/webhook-calm:v2" .

docker push littlebull/webhook-calm:v2

operation = DELETE, webhook不能patch操作。否则报错
Error from server (InternalError): error when deleting "kata-busybox.yaml": Internal error occurred: Internal error occurred: admission webhook "webhook-calm-server.default.svc" attempted to modify the object, which is not supported for this operation

把service停掉，就可以防止apiserver调用webhook了