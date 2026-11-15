#!/bin/bash
user=${1:-username}
passwd=${2:-password}
iface=${3:-ens18}
worker_nodes=(${4:-"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5"})
sudo apt-get install ca-certificates curl gnupg sshpass net-tools -y
host=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)
sudo curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable=traefik --flannel-iface=$iface --node-external-ip $host --kube-apiserver-arg enable-admission-plugins=PodNodeSelector,PodTolerationRestriction -v=1 --log=/var/test-k3s.log" sh -
sleep 30s
sudo cp /etc/rancher/k3s/k3s.yaml /root/.kube/config
echo "export KUBECONFIG=/root/.kube/config" | sudo tee -a /root/.profile >> /dev/null
export token=$(sudo cat /var/lib/rancher/k3s/server/node-token)
for i in ${worker_nodes[@]}; do echo $passwd | sshpass -p $passwd ssh $user@$i -tt -o StrictHostKeyChecking=no "sudo apt-get install ca-certificates curl gnupg -y && curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"agent --server https://$host:6443 --token $token --flannel-iface=$iface --node-external-ip $i\" sh -"; done;
host_name=$(hostname)
sed -i "s/kubernetes.io\/hostname: macropod-1/kubernetes.io\/hostname: $host_name/g" macropod.yaml
sed -i "s/value: \"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5\"/value: \"$host $4\"/g" macropod.yaml
sudo kubectl apply -f macropod.yaml
sed -i "s/value: \"$host $4\"/value: \"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5\"/g" macropod.yaml
sed -i "s/kubernetes.io\/hostname: $host_name/kubernetes.io\/hostname: macropod-1/g" macropod.yaml
# delete if social-network is not used
sleep 60s
kubectl apply -f ./iRoute/deploy/sn-db/namespaces.yml
kubectl apply -f ./iRoute/deploy/sn-db/yaml/
sleep 60s
curl -X POST -d @workflow-definitions/sn-db-op.json http://10.43.190.1:9001/workflow/create 2>/dev/null # change to what is in the macropod.yaml
curl -X POST -d 1 http://10.43.190.1:9001/workflow/invoke/sn-db-op 2>/dev/null
curl -X POST http://10.43.190.1:9001/workflow/delete/sn-db-op 2>/dev/null
#
