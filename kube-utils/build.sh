#!/bin/bash

go build -mod=vendor -gcflags 'all=-N -l' -o kubeutil main.go