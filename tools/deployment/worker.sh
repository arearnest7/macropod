#!/bin/bash
MASTER_IP=$1
TOKEN=$2
curl -sfL https://get.k3s.io | K3S_URL=https://$MASTER_IP:6443 K3S_TOKEN=$TOKEN sh -
mkdir metrics
chmod 777 metrics
cd metrics
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz -O go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.sum -O go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.mod -O go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/metrics.go -O metrics.go
chmod 777 go.sum
chmod 777 go.mod
chmod 777 metrics.go
