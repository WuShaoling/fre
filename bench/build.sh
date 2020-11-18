#!/bin/bash

# build fre_bench
docker run -it --rm -v "$PWD":/go/src golang:1.14 \
  bash -c "cd /go/src && go build -o run run.go && go build -o overlay overlay.go"

scp run overlay root@server:/root/fre/container

rm -f overlay run
