#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get pipelined-checksum-partial > pipelined-checksum-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-checksum-partial
redis-cli -h $REDIS -a $PASSWORD get pipelined-encrypt-partial > pipelined-encrypt-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-encrypt-partial
redis-cli -h $REDIS -a $PASSWORD get pipelined-main-partial > pipelined-main-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-main-partial
redis-cli -h $REDIS -a $PASSWORD get pipelined-zip-partial > pipelined-zip-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-zip-partial
