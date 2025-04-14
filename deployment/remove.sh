#!/bin/bash
worker_nodes=($1)
/usr/local/bin/k3s-uninstall.sh
for i in ${worker_nodes[@]}; do ssh root@$i "/usr/local/bin/k3s-agent-uninstall.sh"; done;
