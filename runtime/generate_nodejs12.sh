#!/bin/bash

# generate rootfs
containerId=$(docker run -d python:3.7)
mkdir rootfs &&
  (docker export "${containerId}" | tar -C rootfs -xvf -) &&
  docker rm -f "${containerId}"

# set dns
cat >rootfs/etc/resolv.conf <<EOF
nameserver 223.5.5.5
nameserver 223.6.6.6
nameserver 8.8.8.8
EOF

# copy test lambdas
mkdir python3.7/code &&
  cp -r code_registry/python3.7/* python3.7/code

### install packages
#docker run -it --rm --privileged -v $PWD:/root python:3.7 bash
#chroot /root/rootfs bash
#pip3 install -i http://pypi.douban.com/simple --trusted-host pypi.douban.com scipy numpy pandas django matplotlib
