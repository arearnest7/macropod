#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get wage-avg > wage-avg-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-avg
redis-cli -h $REDIS -a $PASSWORD get wage-format > wage-format-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-format
redis-cli -h $REDIS -a $PASSWORD get wage-merit > wage-merit-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-merit
redis-cli -h $REDIS -a $PASSWORD get wage-stats > wage-stats-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-stats
redis-cli -h $REDIS -a $PASSWORD get wage-sum > wage-sum-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-sum
redis-cli -h $REDIS -a $PASSWORD get wage-validator > wage-validator-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-validator
redis-cli -h $REDIS -a $PASSWORD get wage-write-merit > wage-write-merit-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-write-merit
redis-cli -h $REDIS -a $PASSWORD get wage-write-raw > wage-write-raw-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del wage-write-raw
