#!/bin/bash

ETCDCTL_CA_FILE=/etc/cni/net.d/calico-tls/etcd-ca ETCDCTL_KEY_FILE=/etc/cni/net.d/calico-tls/etcd-key ETCDCTL_CERT_FILE=/etc/cni/net.d/calico-tls/etcd-cert ETCD_HOST=192.168.6.135 ETCD_PORT=2379 node server.js