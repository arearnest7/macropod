#!/bin/bash
user=$1
interface=$2
metric_nodes=($3)
for i in ${metric_nodes[@]}; do ssh $user@$i "cd /home/$user/metrics && nohup go run metrics.go $interface metrics.csv &"; done;
