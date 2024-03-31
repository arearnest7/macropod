#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get pipelined-checksum > pipelined-checksum-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-checksum
redis-cli -h $REDIS -a $PASSWORD get pipelined-encrypt > pipelined-encrypt-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-encrypt
redis-cli -h $REDIS -a $PASSWORD get pipelined-main > pipelined-main-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-main
redis-cli -h $REDIS -a $PASSWORD get pipelined-zip > pipelined-zip-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-zip
