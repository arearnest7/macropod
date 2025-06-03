#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-payloads/election.json}
N=${4:-1000}
curl -d @workflow-definitions/$BENCHMARK.json -X POST -H "Content-Type: application/json" http://10.43.190.1/create/$BENCHMARK
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1/invoke/$BENCHMARK >> cold-start
number_of_non_running_pods=$(kubectl get pods -n macropod-functions --field-selector=status.phase!=Running --output name | wc -l)
while [ "$number_of_non_running_pods" -ne 0 ]; do
    number_of_non_running_pods=$(kubectl get pods -n macropod-functions --field-selector=status.phase!=Running --output name | wc -l)
    sleep 10
done
sleep 120s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1/invoke/$BENCHMARK >> macropod-$BENCHMARK.csv
curl http://10.43.190.1/delete/$BENCHMARK
number_of_pods=$(kubectl get pods --output name -n macropod-functions | wc -l)
while [ "$number_of_pods" -ne 0 ]; do
    number_of_pods=$(kubectl get pods --output name -n macropod-functions| wc -l)
    sleep 10
done
sleep 180s
