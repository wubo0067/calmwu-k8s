#!/bin/bash

ip=`ifconfig eth0|sed -n 2p|awk  '{ print $2 }'|tr -d 'addr:'`
curl http://${ip}:30001/ping