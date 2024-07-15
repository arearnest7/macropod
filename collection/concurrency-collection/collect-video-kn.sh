#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST full-video video-full $C ../../payloads/video.json $N
mv kn-full-video.csv kn-full-video-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST original-video video-streaming $C ../../payloads/video.json $N
mv kn-original-video.csv kn-original-video-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
./collect-kn.sh $HOST partial-video video-streaming-d $C ../../payloads/video.json $N
mv kn-partial-video.csv kn-partial-video-$C.csv
date -u '+%F %H:%M:%S.%6N %Z'
