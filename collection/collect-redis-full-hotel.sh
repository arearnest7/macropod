#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get hotel-full > hotel-full-log.csv
redis-cli -h $REDIS -a $PASSWORD del hotel-full
