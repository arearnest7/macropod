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
sudo ./deploy.sh $user $passwd $iface "$worker_nodes" 1>/dev/null 2>/dev/null
#sudo ./evaluate.sh 10.43.190.1 $DIR/iRoute 1>/dev/null 2>/dev/null
sudo ./remove.sh $user $passwd "$worker_nodes" 1>/dev/null 2>/dev/null
date
sleep 60s

# Knative
echo "Knative"
date
cd ~/macropod/knative
sudo ./deploy.sh $user $passwd $iface "$worker_nodes" 1>/dev/null 2>/dev/null
#sudo ./evaluate.sh 10.43.190.1 $host $DIR/Knative 1>/dev/null 2>/dev/null
sudo ./remove.sh $user $passwd "$worker_nodes" 1>/dev/null 2>/dev/null
date
sleep 60s

# MacroPod
echo "MacroPod"
date
cd ~/macropod
sudo ./deploy.sh $user $passwd $iface "$worker_nodes" 1>/dev/null 2>/dev/null
#sudo ./evaluate.sh 10.43.190.1 $DIR/MacroPod 1>/dev/null 2>/dev/null
sudo ./remove.sh $user $passwd "$worker_nodes" 1>/dev/null 2>/dev/null
date
