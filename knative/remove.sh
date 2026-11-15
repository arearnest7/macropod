#!/bin/bash
user=${1:-user}
passwd=${2:-password}
worker_nodes=(${3:-"192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5"})
sudo /usr/local/bin/k3s-uninstall.sh
for i in ${worker_nodes[@]}; do sshpass -p $passwd ssh $user@$i -o StrictHostKeyChecking=no "echo $passwd | sudo -S /usr/local/bin/k3s-agent-uninstall.sh"; done;
