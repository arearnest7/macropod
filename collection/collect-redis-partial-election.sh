#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get election-gateway-vevp > election-gateway-vevp-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-gateway-vevp
redis-cli -h $REDIS -a $PASSWORD get election-get-results-partial > election-get-results-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del election-get-results-partial
