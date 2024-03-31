#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get video-recog-partial > video-recog-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del video-recog-partial
redis-cli -h $REDIS -a $PASSWORD get video-streaming-d > video-streaming-d-log.csv
redis-cli -h $REDIS -a $PASSWORD del video-streaming-d
