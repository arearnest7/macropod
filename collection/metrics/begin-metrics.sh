#!/bin/bash
user=$1
interface=$2
metric_nodes=($3)
for i in ${metric_nodes[@]}; do ssh $user@$i "nohup sudo -S env PATH=$PATH:/usr/local/go/bin go run /home/$user/metrics/metrics.go $interface /home/$user/metrics/metrics.csv &"; done;
