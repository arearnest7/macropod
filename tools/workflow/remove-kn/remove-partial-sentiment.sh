#!/bin/bash
kn func delete -p ../../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-db-s
kn func delete -p ../../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-main-rcposc
kn func delete -p ../../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-product-sentiment-prs
kn func delete -p ../../../benchmarks/$1/partial-reduced/sentiment-analysis/sentiment-service-sentiment-srs

