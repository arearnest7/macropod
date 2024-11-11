#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-feature feature-full $C ../../payloads/feature.json $N
mv kn-full-feature.csv kn-full-feature-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-feature feature-orchestrator $C ../../payloads/feature.json $N
mv kn-original-feature.csv kn-original-feature-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-feature feature-orchestrator-wsr $C ../../payloads/feature.json $N
mv kn-partial-feature.csv kn-partial-feature-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-macropod.sh feature $C ../../payloads/feature.json $N
mv macropod-feature.csv macropod-feature-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
