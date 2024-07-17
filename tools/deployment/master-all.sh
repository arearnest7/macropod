#!/bin/bash
host=$1
user=$2
worker_nodes=($3)
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo apt install docker.io
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --kube-apiserver-arg enable-admission-plugins=PodNodeSelector -v=10 --log=/var/test-k3s.log --kube-scheduler-arg=v=10" sh -
sleep 30s
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | sudo tee -a /root/.profile >> /dev/null
for i in ${worker_nodes[@]}; do ssh $user@$i "wget -P /home/$user/ https://raw.githubusercontent.com/arearnest7/macropod/main/tools/deployment/worker.sh -O worker.sh && chmod +x /home/$user/worker.sh && sudo -S /home/$user/worker.sh $host $token"; done;
sudo kubectl apply -f knative/serving-crds.yaml
sudo kubectl apply -f knative/serving-core.yaml
sudo kubectl apply -f knative/istio.yaml
sudo kubectl apply -f knative/net-istio.yaml
sudo kubectl apply -f knative/serving-default-domain.yaml
sudo kubectl apply -f knative/serving-hpa.yaml
sudo kubectl apply -f autoscaler.yaml
sudo kubectl apply -f macropod.yaml
wget https://github.com/knative/client/releases/download/knative-v1.11.2/kn-linux-amd64
mv kn-linux-amd64 kn
chmod +x kn
sudo mv kn /usr/local/bin
wget https://github.com/knative/func/releases/download/knative-v1.12.0/func_linux_amd64
mv func_linux_amd64 kn-func
chmod +x kn-func
sudo mv kn-func /usr/local/bin
sudo apt install hey
mkdir ~/metrics
cd ~/metrics
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz -O go1.22.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /root/.profile >> /dev/null
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.sum -O go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.mod -O go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/metrics.go -O metrics.go
