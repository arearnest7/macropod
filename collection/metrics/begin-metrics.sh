#!/bin/bash
user=$1
interface=$2
metric_nodes=($3)
for i in ${metric_nodes[@]}; do ssh $user@$i "cd /home/$user/metrics; nohup env PATH=$PATH:/usr/local/go/bin go run metrics.go $interface metrics.csv 1>/dev/null 2>/dev/null &"; done;
