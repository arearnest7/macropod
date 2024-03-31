#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get election-gateway > election-gateway-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-gateway
redis-cli -h $REDIS -a $PASSWORD get election-get-results > election-get-results-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-get-results
redis-cli -h $REDIS -a $PASSWORD get election-vote-enqueuer > election-vote-enqueuer-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-vote-enqueuer
redis-cli -h $REDIS -a $PASSWORD get election-vote-processor > election-vote-processor-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-vote-processor
