#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-db-s
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-main-rcposc
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-product-sentiment-prs
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-service-sentiment-srs

