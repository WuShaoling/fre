#!/bin/bash

# build fre_bench
docker run -it -v "$PWD":/go/src golang:1.14 bash -c "cd /go/src && go build -o run run.go"
scp run root@server:/root/fre/container