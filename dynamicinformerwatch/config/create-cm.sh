#!/bin/bash

kubectl delete configmap dynamic-informer-watchres-config -n kata-ns
kubectl create configmap dynamic-informer-watchres-config --from-file=./config.json -n kata-ns