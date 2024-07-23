#!/bin/bash
user=$1
metric_nodes=($2)
for i in ${metric_nodes[@]}; do ssh $user@$i "pkill metrics"; done;
for i in ${metric_nodes[@]}; do scp $user@$i:/home/$user/metrics/metrics.csv $i-metrics.csv; done;
