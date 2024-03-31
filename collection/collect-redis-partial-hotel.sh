#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get hotel-frontend-spgr > hotel-frontend-spgr-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-frontend-spgr
redis-cli -h $REDIS -a $PASSWORD get hotel-recommend-partial > hotel-recommend-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-recommend-partial
redis-cli -h $REDIS -a $PASSWORD get hotel-reserve-partial > hotel-reserve-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-reserve-partial
redis-cli -h $REDIS -a $PASSWORD get hotel-user-partial > hotel-user-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-user-partial
