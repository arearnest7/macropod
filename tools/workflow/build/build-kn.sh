#!/bin/bash
sudo kn func build --push=true --path ../../../benchmarks/kn/original/micro/micro-rpc-a
sudo kn func build --push=true --path ../../../benchmarks/kn/original/micro/micro-rpc-b
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/micro/micro-rpc-a-b
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/feature-generation/feature-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/hotel-app/hotel-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/pipelined-processing/pipelined-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/sentiment-analysis/sentiment-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/serverless-election/election-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/video-analytics/video-full
sudo kn func build --push=true --path ../../../benchmarks/kn/full-reduced/wage-pay/wage-full
sudo kn func build --push=true --path ../../../benchmarks/kn/original/feature-generation/feature-extractor
sudo kn func build --push=true --path ../../../benchmarks/kn/original/feature-generation/feature-orchestrator
sudo kn func build --push=true --path ../../../benchmarks/kn/original/feature-generation/feature-reducer
sudo kn func build --push=true --path ../../../benchmarks/kn/original/feature-generation/feature-status
sudo kn func build --push=true --path ../../../benchmarks/kn/original/feature-generation/feature-wait
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-frontend
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-geo
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-profile
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-rate
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-recommend
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-reserve
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-search
sudo kn func build --push=true --path ../../../benchmarks/kn/original/hotel-app/hotel-user
sudo kn func build --push=true --path ../../../benchmarks/kn/original/pipelined-processing/pipelined-checksum
sudo kn func build --push=true --path ../../../benchmarks/kn/original/pipelined-processing/pipelined-encrypt
sudo kn func build --push=true --path ../../../benchmarks/kn/original/pipelined-processing/pipelined-main
sudo kn func build --push=true --path ../../../benchmarks/kn/original/pipelined-processing/pipelined-zip
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-cfail
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-db
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-main
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-product-or-service
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-product-result
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-product-sentiment
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-read-csv
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-service-result
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-service-sentiment
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-sfail
sudo kn func build --push=true --path ../../../benchmarks/kn/original/sentiment-analysis/sentiment-sns
sudo kn func build --push=true --path ../../../benchmarks/kn/original/serverless-election/election-gateway
sudo kn func build --push=true --path ../../../benchmarks/kn/original/serverless-election/election-get-results
sudo kn func build --push=true --path ../../../benchmarks/kn/original/serverless-election/election-vote-enqueuer
sudo kn func build --push=true --path ../../../benchmarks/kn/original/serverless-election/election-vote-processor
sudo kn func build --push=true --path ../../../benchmarks/kn/original/video-analytics/video-streaming
sudo kn func build --push=true --path ../../../benchmarks/kn/original/video-analytics/video-decoder
sudo kn func build --push=true --path ../../../benchmarks/kn/original/video-analytics/video-recog
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-avg
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-format
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-merit
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-stats
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-sum
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-validator
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-write-merit
sudo kn func build --push=true --path ../../../benchmarks/kn/original/wage-pay/wage-write-raw
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/feature-generation/feature-extractor-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/feature-generation/feature-orchestrator-wsr
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/hotel-app/hotel-frontend-spgr
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/hotel-app/hotel-recommend-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/hotel-app/hotel-reserve-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/hotel-app/hotel-user-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/pipelined-processing/pipelined-checksum-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/pipelined-processing/pipelined-encrypt-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/pipelined-processing/pipelined-main-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/pipelined-processing/pipelined-zip-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/sentiment-analysis/sentiment-db-s
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/sentiment-analysis/sentiment-main-rcposc
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/sentiment-analysis/sentiment-product-sentiment-prs
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/sentiment-analysis/sentiment-service-sentiment-srs
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/serverless-election/election-gateway-vevp
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/serverless-election/election-get-results-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/video-analytics/video-streaming-d
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/video-analytics/video-recog-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/wage-pay/wage-stats-partial
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/wage-pay/wage-sum-amw
sudo kn func build --push=true --path ../../../benchmarks/kn/partial-reduced/wage-pay/wage-validator-fw
