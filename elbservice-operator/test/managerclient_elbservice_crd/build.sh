#!/bin/bash

go build -v -x -mod=vendor -gcflags 'all=-N -l' access.go