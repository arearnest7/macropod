#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get video-decoder > video-decoder-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del video-decoder
redis-cli -h $REDIS -a $PASSWORD get video-recog > video-recog-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del video-recog
redis-cli -h $REDIS -a $PASSWORD get video-streaming > video-streaming-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del video-streaming
