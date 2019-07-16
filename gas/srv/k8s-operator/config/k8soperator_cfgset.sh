#!/bin/bash

consul kv put config/srv/k8soperator "$(cat ./config.json)"