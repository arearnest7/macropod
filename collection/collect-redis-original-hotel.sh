#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get hotel-frontend > hotel-frontend-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-frontend
redis-cli -h $REDIS -a $PASSWORD get hotel-geo > hotel-geo-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-geo
redis-cli -h $REDIS -a $PASSWORD get hotel-profile > hotel-profile-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-profile
redis-cli -h $REDIS -a $PASSWORD get hotel-rate > hotel-rate-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-rate
redis-cli -h $REDIS -a $PASSWORD get hotel-recommend > hotel-recommend-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-recommend
redis-cli -h $REDIS -a $PASSWORD get hotel-reserve > hotel-reserve-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-reserve
redis-cli -h $REDIS -a $PASSWORD get hotel-search > hotel-search-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-search
redis-cli -h $REDIS -a $PASSWORD get hotel-user > hotel-user-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-user
