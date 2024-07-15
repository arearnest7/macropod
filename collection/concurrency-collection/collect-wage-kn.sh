#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-wage wage-full $C ../../payloads/wage.json $N
mv kn-full-wage.csv kn-full-wage-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-wage wage-validator $C ../../payloads/wage.json $N
mv kn-original-wage.csv kn-original-wage-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-wage wage-validator-fw $C ../../payloads/wage.json $N
mv kn-partial-wage.csv kn-partial-wage-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
