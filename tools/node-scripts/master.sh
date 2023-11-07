#!/bin/bash
wget https://github.com/etcd-io/etcd/releases/download/v3.4.16/etcd-v3.4.16-linux-amd64.tar.gz
tar -xvf etcd-v3.4.16-linux-amd64.tar.gz
sudo mv etcd-v3.4.16-linux-amd64 /usr/local/bin/etcd
rm -r etcd*
wget https://github.com/containerd/containerd/releases/download/v1.7.8/containerd-1.7.8-linux-amd64.tar.gz
tar -xvf containerd-1.7.8-linux-amd64.tar.gz
mv containerd-1.7.8-linux-amd64/bin/* /usr/local/bin
rm -r containerd*
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
# Add the repository to Apt sources:
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
curl -sfL https://get.k3s.io | sh -
sudo k3s kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.3/manifests/tigera-operator.yaml
sudo k3s kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.3/manifests/custom-resources.yaml
sudo k3s kubectl delete service -n kube-system traefik
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-crds.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-core.yaml
sudo k3s kubectl apply -l knative.dev/crd-install=true -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/net-istio/releases/download/knative-v1.12.0/net-istio.yaml
sudo k3s kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-default-domain.yaml
wget https://github.com/knative/client/releases/download/knative-v1.11.2/kn-linux-amd64
mv kn-linux-amd64 kn
chmod +x kn
mv kn /usr/local/bin
wget https://github.com/knative/func/releases/download/knative-v1.12.0/func_linux_amd64
mv func_linux_amd64 kn-func
chmod +x kn-func
mv kn-func /usr/local/bin
