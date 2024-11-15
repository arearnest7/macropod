#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-../payloads/election.json}
N=${4:-10000}
curl -d @../../workflow-definitions/$BENCHMARK.json -X POST -H "Content-Type: application/json" http://10.43.190.1/create/$BENCHMARK
sleep 180s
hey -n $C -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1/invoke/$BENCHMARK
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1/invoke/$BENCHMARK >> macropod-$BENCHMARK.csv
logs=$(sudo kubectl get pods -n macropod-functions --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do sudo kubectl logs -n macropod-functions -c user-container $i >> ../../../collection/concurrency-collection/macropod-$SCRIPT-$i.csv; done;
curl http://10.43.190.1/delete/$BENCHMARK
sleep 180s
