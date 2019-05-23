#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

#ROOT=$(cd $(dirname $0)/../../; pwd)

#CA_BUNDLE=$(kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n')
CA_BUNDLE=$(cat key/ca.crt|base64|tr -d "\n")
sed -i "s#CA_BUNDLE#$CA_BUNDLE#g" webhook-calm-deployment.yaml