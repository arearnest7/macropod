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
curl -sSL https://cli.openfaas.com | sudo -E sh
kubectl apply -f ./deploy/OpenFaaS/namespaces.yml
kubectl -n openfaas create secret generic basic-auth --from-literal=basic-auth-user=admin --from-literal=basic-auth-password="$passwd"
kubectl apply -f ./deploy/OpenFaaS/yaml/
sleep 60s
export OPENFAAS_URL=http://127.0.0.1:31112
echo $passwd | faas-cli login --username admin --password-stdin
sleep 60s
kubectl apply -f ./deploy/sn-db/namespaces.yml
kubectl apply -f ./deploy/sn-db/yaml/
faas-cli deploy -f scripts/yamls/sn-db-op.yml 2>/dev/null
sleep 60s
curl -d 1 http://127.0.0.1:31112/function/sn-db-op 2>/dev/null
faas-cli delete -f scripts/yamls/sn-db-op.yml 2>/dev/null
sleep 60s
kubectl apply -f ./deploy/ZooKeeper/local-path-storage.yaml
kubectl apply -f ./deploy/ZooKeeper/zookeeper-config.yaml
kubectl apply -f ./deploy/ZooKeeper/zookeeper-headless-svc.yaml
kubectl apply -f ./deploy/ZooKeeper/zookeeper-statefulset.yaml
sleep 60s
sed -i "s/kubernetes.io\/hostname: macropod-1/kubernetes.io\/hostname: $host_name/g" ../macropod.yaml
sed -i "s/value: \"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5\"/value: \"$host $4\"/g" ../macropod.yaml
sudo kubectl apply -f ../macropod.yaml
sed -i "s/value: \"$host $4\"/value: \"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5\"/g" ../macropod.yaml
sed -i "s/kubernetes.io\/hostname: $host_name/kubernetes.io\/hostname: macropod-1/g" ../macropod.yaml
