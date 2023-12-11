#!/bin/bash
BUILD=${1:-false}
PUSH=${2:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/sentiment-analysis/sentiment-db-s-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/sentiment-analysis/sentiment-main-rcposc-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/sentiment-analysis/sentiment-product-sentiment-prs-wob
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks-wob/partial-reduced/sentiment-analysis/sentiment-service-sentiment-srs-wob

