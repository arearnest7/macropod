#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-sentiment sentiment-full $C ../../payloads/sentiment.json $N
mv kn-full-sentiment.csv kn-full-sentiment-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-sentiment sentiment-main $C ../../payloads/sentiment.json $N
mv kn-original-sentiment.csv kn-original-sentiment-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-sentiment sentiment-main-rcposc $C ../../payloads/sentiment.json $N
mv kn-partial-sentiment.csv kn-partial-sentiment-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
#./collect-macropod.sh sentiment $C ../../payloads/sentiment.json $N
#mv macropod-sentiment.csv macropod-sentiment-$C.csv
#date -u '+%F %H:%M:%S.%6N %Z'
