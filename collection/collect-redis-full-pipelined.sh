#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get pipelined-full > pipelined-full-log.csv
redis-cli -h $REDIS -a $PASSWORD del pipelined-full
