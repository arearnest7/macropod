#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get sentiment-db-s > sentiment-db-s-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-db-s
redis-cli -h $REDIS -a $PASSWORD get sentiment-main-rcposc > sentiment-main-rcposc-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-main-rcposc
redis-cli -h $REDIS -a $PASSWORD get sentiment-product-prs > sentiment-product-prs-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-product-prs
redis-cli -h $REDIS -a $PASSWORD get sentiment-service-srs > sentiment-service-srs-log.csv
redis-cli -h $REDIS -a $PASSWORD del sentiment-service-srs
