#!/bin/bash

docker build -t littlebull/testnetlink:v1 .
docker run -itd littlebull/testnetlink:v1
container=$(docker ps | awk 'NR > 1 {print $1; exit}')
echo "-----"
echo $container
echo "-----"
docker cp $container":/netlink/testnetlink" ./
docker stop $container