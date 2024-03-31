#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get feature-extractor > feature-extractor-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-extractor
redis-cli -h $REDIS -a $PASSWORD get feature-orchestrator > feature-orchestrator-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-orchestrator
redis-cli -h $REDIS -a $PASSWORD get feature-reducer > feature-reducer-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-reducer
redis-cli -h $REDIS -a $PASSWORD get feature-status > feature-status-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-status
redis-cli -h $REDIS -a $PASSWORD get feature-wait > feature-wait-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-wait
