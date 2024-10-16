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
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://$ENTRY.knative-functions.$HOST.sslip.io >> ../../../collection/fixed-collection/kn-$SCRIPT.csv
logs=$(sudo kubectl get pods -n knative-functions --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do sudo kubectl logs -n knative-functions -c user-container $i >> ../../../collection/fixed-collection/kn-$SCRIPT-$i.csv; done;
../remove-kn/remove-$SCRIPT.sh
cd ../../../collection/fixed-collection/
sleep 180s
