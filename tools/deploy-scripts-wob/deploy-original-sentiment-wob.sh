#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-cfail-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-db-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-main-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-product-or-service-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-product-result-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-product-sentiment-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-read-csv-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-service-result-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-service-sentiment-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-sfail-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/original/sentiment-analysis/sentiment-sns-wob

