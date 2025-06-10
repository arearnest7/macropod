#!/bin/bash
user=$1
iface=$2
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
for i in ${worker_nodes[@]}; do ssh $user@$i -tt "sudo curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"agent --server https://$host:6443 --token $token --flannel-iface=$iface --node-external-ip $i\" sh -"; done;
host_name=$(hostname)
sudo kubectl taint nodes $host_name master-node=master-node:NoSchedule
sudo kubectl apply -f deployment/knative/serving-crds.yaml
sudo kubectl apply -f deployment/knative/serving-core.yaml
sudo kubectl apply -f deployment/knative/istio-ns.yaml
sudo kubectl apply -f deployment/knative/istio.yaml
sudo kubectl apply -f deployment/knative/net-istio.yaml
sudo kubectl apply -f deployment/knative/autoscaler.yaml
sleep 60s
sudo kubectl get daemonset -A -o jsonpath='{range .items[*]}{.metadata.name}{" -n "}{.metadata.namespace}{"\n"}{end}' | while read -r line; do sudo kubectl patch daemonset $line --patch-file deployment/knative/daemonset.yaml; done;
sudo kubectl apply -f deployment/knative/serving-default-domain.yaml
sudo kubectl apply -f deployment/macropod.yaml
wget https://github.com/knative/client/releases/download/knative-v1.11.2/kn-linux-amd64
mv kn-linux-amd64 kn
chmod +x kn
sudo mv kn /usr/local/bin
wget https://github.com/knative/func/releases/download/knative-v1.12.0/func_linux_amd64
mv func_linux_amd64 kn-func
chmod +x kn-func
sudo mv kn-func /usr/local/bin
