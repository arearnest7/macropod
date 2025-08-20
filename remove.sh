#!/bin/bash
user=${1:-user}
passwd=${2:-password}
worker_nodes=(${3:-"192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24"})
sudo /usr/local/bin/k3s-uninstall.sh
for i in ${worker_nodes[@]}; do sshpass -p $passwd ssh $user@$i -o StrictHostKeyChecking=no "echo $passwd | sudo -S /usr/local/bin/k3s-agent-uninstall.sh"; done;
