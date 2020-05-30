#!/bin/bash

kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/crds/k8s.calmwu.org_elbservices_crd.yaml
kubectl delete cm elbservice-cm -n calmwu-namespace
kubectl create cm elbservice-cm --from-file=deploy/config/elbservice_config.json -n calmwu-namespace
