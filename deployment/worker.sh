#!/bin/bash
MASTER_IP=$1
iface=$2
TOKEN=$3
sudo apt-get install ca-certificates curl gnupg
worker=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="agent --server https://$MASTER_IP:6443 --token $TOKEN --flannel-iface $iface --node-external-ip $worker" sh -
mkdir metrics
chmod 777 metrics
cd metrics
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz -O go.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.sum -O go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.mod -O go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/metrics.go -O metrics.go
chmod 777 go.sum
chmod 777 go.mod
chmod 777 metrics.go
