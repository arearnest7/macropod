#!/bin/bash
HOST=${1:-localhost}
SCRIPT=${2:-full-election}
TYPE=${3:-kn}
ENTRY=${4:-election-full}
C=${5:-1}
PAYLOAD=${6:-../payloads/election.json}
N=${7:-10000}
cd ../../tools/deploy-functions-kn/
./deploy-$SCRIPT.sh $TYPE
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://$ENTRY.default.$HOST.sslip.io >> ../../collection/benchmark-collection/kn-$SCRIPT.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> ../../collection/benchmark-collection/kn-$SCRIPT-$i.csv; done;
../remove-functions-kn/remove-$SCRIPT.sh $TYPE
cd ../../collection/benchmark-collection/
sleep 180s
