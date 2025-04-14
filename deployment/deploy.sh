#!/bin/bash
host=${1:-127.0.0.1}
worker_nodes=($2)
apt-get update
apt-get install ca-certificates curl gnupg -y
apt install docker.io -y
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --node-external-ip $host --flannel-external-ip --kube-apiserver-arg enable-admission-plugins=PodNodeSelector,PodTolerationRestriction -v=1 --log=/var/test-k3s.log" sh -
sleep 30s
cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | tee -a /root/.profile >> /dev/null
export token=$(cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do ssh root@$i "apt-get install ca-certificates curl gnupg -y && curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"agent --server https://$host:6443 --token $token --node-external-ip $i\" sh -"; done;
host_name=$(hostname)
kubectl taint nodes $host_name master-node=master-node:NoSchedule
kubectl apply -f macropod.yaml
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz -O go.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
