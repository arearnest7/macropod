#!/bin/bash
HOST=${1:-127.0.0.1}
SCRIPT=${2:-full-election}
ENTRY=${3:-election-full}
C=${4:-1}
PAYLOAD=${5:-payloads/election.json}
N=${6:-1000}
cd kn-scripts/
./deploy-$SCRIPT.sh
cd ../
sleep 60s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://$ENTRY.knative-functions.$HOST.sslip.io >> cold-start
number_of_non_running_pods=$(kubectl get pods -n knative-functions --field-selector=status.phase!=Running --output name | wc -l)
while [ "$number_of_non_running_pods" -ne 0 ]; do
    number_of_non_running_pods=$(kubectl get pods -n knative-functions --field-selector=status.phase!=Running --output name | wc -l)
    sleep 10
done
sleep 120s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://$ENTRY.knative-functions.$HOST.sslip.io >> kn-$SCRIPT.csv
logs=$(sudo kubectl get pods -n knative-functions --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do sudo kubectl logs -n knative-functions -c user-container $i >> kn-$SCRIPT-$i.csv; done;
cd kn-scripts/
./remove-$SCRIPT.sh
cd ../
number_of_pods=$(kubectl get pods --output name -n knative-functions | wc -l)
while [ "$number_of_pods" -ne 0 ]; do
    number_of_pods=$(kubectl get pods --output name -n knative-functions| wc -l)
    sleep 10
done
sleep 180s
