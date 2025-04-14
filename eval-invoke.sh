#!/bin/bash
target=${1:-10.43.190.1}
workflow=${2:-election-agg}
results=${3:-results}

id=$(curl -d @workflow-definitions/$workflow -X POST http://$target/eval/start)
curl http://$target/eval/metrics/$id >> $results/$workflow-metrics.csv
curl http://$target/eval/latency/$id >> $results/$workflow-latency.csv
curl http://$target/eval/summary/$id >> $results/$workflow-summary.csv
