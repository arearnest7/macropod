#!/bin/bash
MASTER_IP=$1
TOKEN=$2
wget https://github.com/etcd-io/etcd/releases/download/v3.4.16/etcd-v3.4.16-linux-amd64.tar.gz
tar -xvf etcd-v3.4.16-linux-amd64.tar.gz
sudo mv etcd-v3.4.16-linux-amd64 /usr/local/bin/etcd
rm -r etcd*
wget https://github.com/containerd/containerd/releases/download/v1.7.8/containerd-1.7.8-linux-amd64.tar.gz
tar -xvf containerd-1.7.8-linux-amd64.tar.gz
mv containerd-1.7.8-linux-amd64/bin/* /usr/local/bin
rm -r containerd*
curl -sfL https://get.k3s.io | K3S_URL=https://$MASTER_IP:6443 K3S_TOKEN=$TOKEN sh -
