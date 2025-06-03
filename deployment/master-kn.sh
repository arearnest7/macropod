#!/bin/bash
iface=$1
user=$2
worker_nodes=($3)
host=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo apt install docker.io
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --flannel-iface=$iface --node-external-ip $host --flannel-external-ip --kube-apiserver-arg enable-admission-plugins=PodNodeSelector,PodTolerationRestriction -v=1 --log=/var/test-k3s.log" sh -
sleep 30s
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | sudo tee -a /root/.profile >> /dev/null
export token=$(sudo cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do ssh $user@$i -tt "wget -P /home/$user/ https://raw.githubusercontent.com/arearnest7/macropod/main/tools/deployment/worker.sh -O worker.sh && chmod +x /home/$user/worker.sh && sudo /home/$user/worker.sh $host $iface $token"; done;
host_name=$(hostname)
kubectl taint nodes $host_name master-node=master-node:NoSchedule
sudo kubectl apply -f knative/serving-crds.yaml
sudo kubectl apply -f knative/serving-core.yaml
sudo kubectl apply -f knative/istio-ns.yaml
sudo kubectl apply -f knative/istio.yaml
sudo kubectl apply -f knative/net-istio.yaml
sudo kubectl apply -f knative/autoscaler.yaml
sleep 60s
sudo kubectl get daemonset -A -o jsonpath='{range .items[*]}{.metadata.name}{" -n "}{.metadata.namespace}{"\n"}{end}' | while read -r line; do sudo kubectl patch daemonset $line --patch-file knative/daemonset.yaml; done
sleep 60s
sudo kubectl apply -f knative/serving-default-domain.yaml
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
chmod 777 ~/metrics
cd ~/metrics
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
