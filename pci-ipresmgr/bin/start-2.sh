#!/bin/bash

ip=`ifconfig eth0|sed -n 2p|awk  '{ print $2 }'|tr -d 'addr:'`
nohup ./amd64/ipresmgr-srv --ip=${ip} --port=30002 --logpath=../log --id=2 --conf=../config/ipresmgr_cfg.json >ipresmgr-srv.out 2>&1 &  