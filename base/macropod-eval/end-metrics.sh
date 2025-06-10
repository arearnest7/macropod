#!/bin/bash
user=$1
results=${2:-results}
metric_nodes=($3)
for i in ${metric_nodes[@]}; do ssh $user@$i "pkill metrics"; done;
for i in ${metric_nodes[@]}; do scp $user@$i:/home/$user/metrics/metrics.csv $results/$i-metrics.csv; done;
