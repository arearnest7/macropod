#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
TYPE=${3:-multi-oci}
redis-cli -h $REDIS -a $PASSWORD get sentiment-cfail > sentiment-cfail-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-cfail
redis-cli -h $REDIS -a $PASSWORD get sentiment-db > sentiment-db-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-db
redis-cli -h $REDIS -a $PASSWORD get sentiment-main > sentiment-main-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-main
redis-cli -h $REDIS -a $PASSWORD get sentiment-product-or-service > sentiment-product-or-service-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-product-or-service
redis-cli -h $REDIS -a $PASSWORD get sentiment-product-result > sentiment-product-result-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-product-result
redis-cli -h $REDIS -a $PASSWORD get sentiment-product-sentiment > sentiment-product-sentiment-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-product-sentiment
redis-cli -h $REDIS -a $PASSWORD get sentiment-read-csv > sentiment-read-csv-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-read-csv
redis-cli -h $REDIS -a $PASSWORD get sentiment-service-result > sentiment-service-result-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-service-result
redis-cli -h $REDIS -a $PASSWORD get sentiment-service-sentiment > sentiment-service-sentiment-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-service-sentiment
redis-cli -h $REDIS -a $PASSWORD get sentiment-sfail > sentiment-sfail-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-sfail
redis-cli -h $REDIS -a $PASSWORD get sentiment-sns > sentiment-sns-$TYPE-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-sns
