#!/bin/bash
PATH=${1:-kn}
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/sentiment-analysis/sentiment-db-s
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/sentiment-analysis/sentiment-main-rcposc
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/sentiment-analysis/sentiment-product-sentiment-prs
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$PATH/partial-reduced/sentiment-analysis/sentiment-service-sentiment-srs

