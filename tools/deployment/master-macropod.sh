#!/bin/bash
host=$1
user=$2
worker_nodes=($3)
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo apt install docker.io
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik" sh -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
mkdir ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chmod 644 ~/.kube/config
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml" >> ~/.profile
export token=$(sudo cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do ssh $user@$i "wget -P /home/$user/ https://raw.githubusercontent.com/arearnest7/macropod/main/tools/deployment/worker.sh && sudo -S /home/$user/worker.sh $host $token"; done;
sudo k3s kubectl apply -f macropod-ingress.yaml
sudo k3s kubectl create clusterrolebinding default-admin --clusterrole=admin --serviceaccount=default:default
sudo apt install hey
mkdir ~/metrics
cd ~/metrics
wget https://go.dev/dl/go1.22.3.src.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/metrics.go