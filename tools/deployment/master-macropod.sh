#!/bin/bash
iface=$1
user=$2
worker_nodes=($3)
host=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo apt install docker.io
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --flannel-iface=$iface --node-external-ip $host --kube-apiserver-arg enable-admission-plugins=PodNodeSelector,PodTolerationRestriction -v=1 --log=/var/test-k3s.log" sh -
sleep 30s
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | sudo tee -a /root/.profile >> /dev/null
export token=$(sudo cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do ssh $user@$i "wget -P /home/$user/ https://raw.githubusercontent.com/arearnest7/macropod/main/tools/deployment/worker.sh -O worker.sh && chmod +x /home/$user/worker.sh && sudo -S /home/$user/worker.sh $host $iface $token"; done;
host_name=$(hostname)
kubectl taint nodes $host_name master-node=master-node:NoSchedule
sudo kubectl apply -f macropod.yaml
sudo apt install hey
mkdir ~/metrics
chmod 777 ~/metrics
cd ~/metrics
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
