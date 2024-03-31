#!/bin/bash
REDIS=${1:-127.0.0.1}
PASSWORD=${2:-password}
redis-cli -h $REDIS -a $PASSWORD get feature-extractor-partial > feature-extractor-partial-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-extractor-partial
redis-cli -h $REDIS -a $PASSWORD get feature-orchestrator-wsr > feature-orchestrator-wsr-log.csv
redis-cli -h $REDIS -a $PASSWORD del feature-orchestrator-wsr
