#!/bin/bash

kubectl get csr my-svc.calm-space -o jsonpath='{.status.certificate}' | base64 --decode > server.crt