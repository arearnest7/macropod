#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-pipelined pipelined-full $C ../../payloads/pipelined.json $N
mv kn-full-pipelined.csv kn-full-pipelined-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-pipelined pipelined-main $C ../../payloads/pipelined.json $N
mv kn-original-pipelined.csv kn-original-pipelined-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-pipelined pipelined-main-partial $C ../../payloads/pipelined.json $N
mv kn-partial-pipelined.csv kn-partial-pipelined-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
#./collect-macropod.sh pipelined $C ../../payloads/pipelined.json $N
#mv macropod-pipelined.csv macropod-pipelined-$C.csv
#date -u '+%F %H:%M:%S.%6N %Z'
