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
for i in ${worker_nodes[@]}; do scp worker.sh $user@$i:/home/$user && ssh $user@$i "sudo -S /home/$user/worker.sh $host $token"; done;
sudo k3s kubectl apply -f macropod-ingress.yaml
sudo k3s kubectl create clusterrolebinding default-admin --clusterrole=admin --serviceaccount=default:default
sudo apt install hey
