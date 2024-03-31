#!/bin/bash
HOST=${1:-localhost}
SCRIPT=${2:-full-election}
TYPE=${3:-kn}
ENTRY=${4:-election-full}
REDIS=${5:-127.0.0.1}
PASSWORD=${6:-password}
C=${7:-1}
PAYLOAD=${8:-../payloads/election.json}
N=${9:-10000}
cd ../tools/deploy-functions-kn/
./deploy-$SCRIPT.sh $TYPE
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://$ENTRY.default.$HOST.sslip.io >> $TYPE-$SCRIPT.csv
mv $TYPE-$SCRIPT.csv ../../collection/
../remove-functions-kn/remove-$SCRIPT.sh $TYPE
cd ../../collection
./collect-redis-$SCRIPT.sh $REDIS $PASSWORD kn
sleep 180s
