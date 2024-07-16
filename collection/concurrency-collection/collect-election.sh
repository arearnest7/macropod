#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-election election-full $C ../../payloads/election.json $N
mv kn-full-election.csv kn-full-election-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-election election-gateway $C ../../payloads/election.json $N
mv kn-original-election.csv kn-original-election-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-election election-gateway-vevp $C ../../payloads/election.json $N
mv kn-partial-election.csv kn-partial-election-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-macropod.sh election $C ../../payloads/election.json $N
mv macropod-election.csv macropod-election-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
