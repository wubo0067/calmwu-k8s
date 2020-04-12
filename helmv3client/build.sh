#!/bin/bash

GO111MODULE=on go build -mod=vendor -gcflags 'all=-N -l' helmv3client.go helmtemplate.go helminstall.go helmstatus.go patchDeploymentTemplate.go patchServiceTemplate.go patchtemplate.go