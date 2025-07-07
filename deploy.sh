#!/bin/bash
user=${1:-username}
passwd=${2:-password}
iface=${3:-enp0s3}
worker_nodes=(${4:-"192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24"})
host=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)
sudo apt-get install ca-certificates curl gnupg sshpass
sudo curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --flannel-iface=$iface --node-external-ip $host --kube-apiserver-arg enable-admission-plugins=PodNodeSelector,PodTolerationRestriction -v=1 --log=/var/test-k3s.log" sh -
sleep 30s
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | sudo tee -a /root/.profile >> /dev/null
export token=$(sudo cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do echo $passwd | sshpass -p $passwd ssh $user@$i -tt -o StrictHostKeyChecking=no "sudo apt-get install ca-certificates curl gnupg -y && curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"agent --server https://$host:6443 --token $token --flannel-iface=$iface --node-external-ip $i\" sh -"; done;
host_name=$(hostname)
sudo kubectl taint nodes $host_name master-node=master-node:NoSchedule
# omit below if knative is not needed
sudo kubectl apply -f deployment/knative/serving-crds.yaml
sudo kubectl apply -f deployment/knative/serving-core.yaml
sudo kubectl apply -f deployment/knative/istio-ns.yaml
sudo kubectl apply -f deployment/knative/istio.yaml
sudo kubectl apply -f deployment/knative/net-istio.yaml
sudo kubectl apply -f deployment/knative/autoscaler.yaml
sleep 60s
sudo kubectl get daemonset -A -o jsonpath='{range .items[*]}{.metadata.name}{" -n "}{.metadata.namespace}{"\n"}{end}' | while read -r line; do sudo kubectl patch daemonset $line --patch-file deployment/knative/daemonset.yaml; done;
sudo kubectl apply -f deployment/knative/serving-default-domain.yaml
sleep 60s
#
sed -i "s/kubernetes.io\/hostname: sysdev-tamu-1/kubernetes.io\/hostname: $host_name/g" deployment/macropod.yaml
sed -i "s/value: \"192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24\"/value: \"$4\"/g" deployment/macropod.yaml
sudo kubectl apply -f deployment/macropod.yaml
sed -i "s/\"value: \"$4\"/value: \"192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24\"/g" deployment/macropod.yaml
sed -i "s/kubernetes.io\/hostname: $host_name/kubernetes.io\/hostname: sysdev-tamu-1/g" deployment/macropod.yaml
