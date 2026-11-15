#!/bin/bash
user=${1:-username}
passwd=${2:-password}
iface=${3:-ens18}
worker_nodes=${4:-"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5"}
DIR=${5:-~/results}

host=$(/sbin/ifconfig $iface | grep -i mask | awk '{print $2}'| cut -f2 -d:)

mkdir $DIR
mkdir $DIR/iRoute
mkdir $DIR/Knative
mkdir $DIR/MacroPod

# iRoute
echo "iRoute"
date
cd ~/macropod/iRoute
sudo ./deploy.sh $user $passwd $iface $worker_nodes
sudo ./evaluate.sh 10.43.190.1 $host $DIR/iRoute
sudo ./remove.sh $user $passwd $worker_nodes
date
sleep 60s

# Knative
echo "Knative"
date
cd ~/macropod/knative
sudo ./deploy.sh $user $passwd $iface $worker_nodes
sudo ./evaluate.sh 10.43.190.1 $DIR/Knative
sudo ./remove.sh $user $passwd $worker_nodes
date
sleep 60s

# MacroPod
echo "MacroPod"
date
cd ~/macropod
sudo ./deploy.sh $user $passwd $iface $worker_nodes
sudo ./evaluate.sh 10.43.190.1 $DIR/MacroPod
sudo ./remove.sh $user $passwd $worker_nodes
date
