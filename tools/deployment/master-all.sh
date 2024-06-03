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
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-crds.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-core.yaml
sudo k3s kubectl apply -l knative.dev/crd-install=true -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/net-istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-default-domain.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.2/serving-hpa.yaml
sudo k3s kubectl apply -f autoscaler.yaml
sudo k3s kubectl apply -f macropod-ingress.yaml
sudo k3s kubectl create clusterrolebinding default-admin --clusterrole=admin --serviceaccount=default:default
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
wget https://go.dev/dl/go1.22.3.src.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.sum
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/go.mod
wget https://raw.githubusercontent.com/arearnest7/macropod/main/tools/collection/metrics/metrics.go
