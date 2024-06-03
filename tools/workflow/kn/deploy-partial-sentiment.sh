#!/bin/bash
BUILD=${2:-false}
PUSH=${3:-false}
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-db-s
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-main-rcposc
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-product-sentiment-prs
kn func deploy --build=$BUILD --push=$PUSH --path ../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-service-sentiment-srs

