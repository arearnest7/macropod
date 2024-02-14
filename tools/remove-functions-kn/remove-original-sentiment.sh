#!/bin/bash
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-cfail
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-db
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-main
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-product-or-service
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-product-result
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-product-sentiment
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-read-csv
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-service-result
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-service-sentiment
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-sfail
kn func delete -p ../../benchmarks/$1/original/sentiment-analysis/sentiment-sns

