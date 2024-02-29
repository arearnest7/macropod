#!/bin/bash
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo apt install docker.io
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik" sh -
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-crds.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-core.yaml
sudo k3s kubectl apply -l knative.dev/crd-install=true -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/net-istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-default-domain.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.2/serving-hpa.yaml
sudo k3s kubectl apply -f autoscaler.yaml
wget https://github.com/knative/client/releases/download/knative-v1.11.2/kn-linux-amd64
mv kn-linux-amd64 kn
chmod +x kn
sudo mv kn /usr/local/bin
rm kn
wget https://github.com/knative/func/releases/download/knative-v1.12.0/func_linux_amd64
mv func_linux_amd64 kn-func
chmod +x kn-func
sudo mv kn-func /usr/local/bin
rm kn-func
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install redis
echo "USE BELOW K3S_TOKEN"
sudo cat /var/lib/rancher/k3s/server/node-token
