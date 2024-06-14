#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-../payloads/election.json}
N=${4:-10000}
curl -d @../../workflow-definitions/$BENCHMARK.json -X POST -H "Content-Type: application/json" http://10.43.190.1/create/$BENCHMARK
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1/invoke/$BENCHMARK >> macropod-$BENCHMARK.csv
curl http://10.43.190.1/logs/$BENCHMARK >> macropod-$BENCHMARK-log-bundle.csv
curl http://10.43.190.1/delete/$BENCHMARK
sleep 180s
