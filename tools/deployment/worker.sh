#!/bin/bash
MASTER_IP=$1
TOKEN=$2
curl -sfL https://get.k3s.io | K3S_URL=https://$MASTER_IP:6443 K3S_TOKEN=$TOKEN sh -
mkdir metrics
cd metrics
wget https://go.dev/dl/go1.22.3.src.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
wget https://raw.githubusercontent.com/arearnest7/macropod/test-branch/tools/collection/metrics/go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/test-branch/tools/collection/metrics/go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/test-branch/tools/collection/metrics/metrics.go
