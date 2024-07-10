#!/bin/bash
HOST=${1:-localhost}
SCRIPT=${2:-full-election}
ENTRY=${3:-election-full}
C=${4:-1}
PAYLOAD=${5:-../payloads/election.json}
N=${6:-10000}
cd ../../tools/workflow/kn/
./deploy-$SCRIPT.sh
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://$ENTRY.default.$HOST.sslip.io >> ../../../collection/concurrency-collection/kn-$SCRIPT.csv
../remove-kn/remove-$SCRIPT.sh
cd ../../../collection/concurrency-collection/
sleep 180s
