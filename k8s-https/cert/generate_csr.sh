#!/bin/bash

cat <<EOF | cfssl_linux-amd64 genkey - | cfssljson_linux-amd64 -bare server
{
  "hosts": [
    "my-svc.calm-space.svc.cluster.local"
  ],
  "CN": "my-svc.calm-space.svc.cluster.local",
  "key": {
    "algo": "ecdsa",
    "size": 256
  },
  "names": [
    {
      "C": "CN",
      "ST": "Wuhan",
      "L": "Wuhan",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF