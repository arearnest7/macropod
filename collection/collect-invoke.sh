#!/bin/bash
HOST=${1:-localhost}
SCRIPT=${2:-full-election}
TYPE=${3:-kn}
ENTRY=${4:-election-full}
C=${5:-1}
PAYLOAD=${6:-../payloads/election.json}
N=${7:-10000}
cd ../tools/deploy-functions-kn/
./deploy-$SCRIPT.sh $TYPE
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://$ENTRY.default.$HOST.sslip.io >> $TYPE-$SCRIPT.csv
mv $TYPE-$SCRIPT.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> $TYPE-$i.csv; done;
../remove-functions-kn/remove-$SCRIPT.sh $TYPE
cd ../../collection
sleep 180s
