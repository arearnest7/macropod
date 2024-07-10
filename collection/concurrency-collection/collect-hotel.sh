#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-hotel hotel-full $C ../../payloads/hotel.json $N
mv kn-full-hotel.csv kn-full-hotel-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-hotel hotel-frontend $C ../../payloads/hotel.json $N
mv kn-original-hotel.csv kn-original-hotel-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-hotel hotel-frontend-spgr $C ../../payloads/hotel.json $N
mv kn-partial-hotel.csv kn-partial-hotel-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
