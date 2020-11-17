#!/bin/bash

# overlayfs 挂载
# mount -t overlay overlay -o lowerdir=rootfs,upperdir=upper,workdir=worker merge
# mount --bind containers rootfs/containers
# mount -t overlay overlay -o lowerdir=rootfs,upperdir=containers/c1/upper,workdir=containers/c1/worker containers/c1/merge

# 启动测试环境
docker run -it --privileged -v $PWD:/root python:3.7 bash
