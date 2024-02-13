#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-cfail
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-db
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-main
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-product-or-service
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-product-result
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-product-sentiment
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-read-csv
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-service-result
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-service-sentiment
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-sfail
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/sentiment-sns

