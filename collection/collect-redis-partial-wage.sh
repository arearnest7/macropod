#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get wage-stats-partial > wage-stats-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-stats-partial
redis-cli -h $REDIS -a $PASSWORD get wage-sum-amw > wage-sum-amw-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-sum-amw
redis-cli -h $REDIS -a $PASSWORD get wage-validator-fw > wage-validator-fw-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-validator-fw
