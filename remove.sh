#!/bin/bash
user=$1
worker_nodes=($2)
sudo /usr/local/bin/k3s-uninstall.sh
for i in ${worker_nodes[@]}; do ssh $user@$i -tt "sudo /usr/local/bin/k3s-agent-uninstall.sh"; done;
