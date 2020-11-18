#!/bin/bash

image=python:3.7
name=python3.7

# generate rootfs
containerId=$(docker run -d ${image})
mkdir ${name} &&
  (docker export "${containerId}" | tar -C ${name} -xvf -) &&
  docker rm -f "${containerId}"

# set dns
cat >${name}/etc/resolv.conf <<EOF
nameserver 223.5.5.5
nameserver 223.6.6.6
nameserver 8.8.8.8
EOF

# 拷贝 demo 代码
mkdir ${name}/code &&
  cp -r code_registry/${name}/* ${name}/code

## 安装 demo 代码的依赖包
docker run -it --rm --privileged -v $PWD:/root python:3.7 bash -c \
  "chroot /root/${name} pip3 install -i http://pypi.douban.com/simple --trusted-host pypi.douban.com scipy numpy pandas django matplotlib"
