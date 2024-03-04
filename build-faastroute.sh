#!/bin/bash
ID=${1:-arearnest7}
docker build ./templates/faastroute-go -t $ID/faastroute-go:latest && docker push $ID/faastroute-go:latest
docker build ./templates/faastroute-node -t $ID/faastroute-node:latest && docker push $ID/faastroute-node:latest
docker build ./templates/faastroute-python -t $ID/faastroute-python:latest && docker push $ID/faastroute-python:latest
docker image prune -a -f
docker build ./benchmarks/micro/micro-rpc-a-faastroute -t $ID/micro-rpc-a:faastroute && docker push $ID/micro-rpc-a:faastroute
docker build ./benchmarks/micro/micro-rpc-b-faastroute -t $ID/micro-rpc-b:faastroute && docker push $ID/micro-rpc-b:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/feature-generation/feature-extractor -t $ID/feature-extractor:faastroute && docker push $ID/feature-extractor:faastroute
docker build ./benchmarks/faastroute/original/feature-generation/feature-orchestrator -t $ID/feature-orchestrator:faastroute && docker push $ID/feature-orchestrator:faastroute
docker build ./benchmarks/faastroute/original/feature-generation/feature-reducer -t $ID/feature-reducer:faastroute && docker push $ID/feature-reducer:faastroute
docker build ./benchmarks/faastroute/original/feature-generation/feature-status -t $ID/feature-status:faastroute && docker push $ID/feature-status:faastroute
docker build ./benchmarks/faastroute/original/feature-generation/feature-wait -t $ID/feature-wait:faastroute && docker push $ID/feature-wait:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:faastroute && docker push $ID/hotel-frontend:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-geo -t $ID/hotel-geo:faastroute && docker push $ID/hotel-geo:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-profile -t $ID/hotel-profile:faastroute && docker push $ID/hotel-profile:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-rate -t $ID/hotel-rate:faastroute && docker push $ID/hotel-rate:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:faastroute && docker push $ID/hotel-recommend:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:faastroute && docker push $ID/hotel-reserve:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-search -t $ID/hotel-search:faastroute && docker push $ID/hotel-search:faastroute
docker build ./benchmarks/faastroute/original/hotel-app/hotel-user -t $ID/hotel-user:faastroute && docker push $ID/hotel-user:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/pipelined-processing/pipelined-checksum -t $ID/pipelined-checksum:faastroute && docker push $ID/pipelined-checksum:faastroute
docker build ./benchmarks/faastroute/original/pipelined-processing/pipelined-encrypt -t $ID/pipelined-encrypt:faastroute && docker push $ID/pipelined-encrypt:faastroute
docker build ./benchmarks/faastroute/original/pipelined-processing/pipelined-main -t $ID/pipelined-main:faastroute && docker push $ID/pipelined-main:faastroute
docker build ./benchmarks/faastroute/original/pipelined-processing/pipelined-zip -t $ID/pipelined-zip:faastroute && docker push $ID/pipelined-zip:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-cfail -t $ID/sentiment-cfail:faastroute && docker push $ID/sentiment-cfail:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-main -t $ID/sentiment-main:faastroute && docker push $ID/sentiment-main:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-product-result -t $ID/sentiment-product-result:faastroute && docker push $ID/sentiment-product-result:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-read-csv -t $ID/sentiment-read-csv:faastroute && docker push $ID/sentiment-read-csv:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-service-sentiment -t $ID/sentiment-service-sentiment:faastroute && docker push $ID/sentiment-service-sentiment:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-sns -t $ID/sentiment-sns:faastroute && docker push $ID/sentiment-sns:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-db -t $ID/sentiment-db:faastroute && docker push $ID/sentiment-db:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-product-or-service -t $ID/sentiment-product-or-service:faastroute && docker push $ID/sentiment-product-or-service:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-product-sentiment -t $ID/sentiment-product-sentiment:faastroute && docker push $ID/sentiment-product-sentiment:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-service-result -t $ID/sentiment-service-result:faastroute && docker push $ID/sentiment-service-result:faastroute
docker build ./benchmarks/faastroute/original/sentiment-analysis/sentiment-sfail -t $ID/sentiment-sfail:faastroute && docker push $ID/sentiment-sfail:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/serverless-election/election-gateway -t $ID/election-gateway:faastroute && docker push $ID/election-gateway:faastroute
docker build ./benchmarks/faastroute/original/serverless-election/election-get-results -t $ID/election-get-results:faastroute && docker push $ID/election-get-results:faastroute
docker build ./benchmarks/faastroute/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:faastroute && docker push $ID/election-vote-enqueuer:faastroute
docker build ./benchmarks/faastroute/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:faastroute && docker push $ID/election-vote-processor:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/video-analytics/video-decoder -t $ID/video-decoder:faastroute && docker push $ID/video-decoder:faastroute
docker build ./benchmarks/faastroute/original/video-analytics/video-recog -t $ID/video-recog:faastroute && docker push $ID/video-recog:faastroute
docker build ./benchmarks/faastroute/original/video-analytics/video-streaming -t $ID/video-streaming:faastroute && docker push $ID/video-streaming:faastroute
docker image prune -a -f
docker build ./benchmarks/faastroute/original/wage-pay/wage-avg -t $ID/wage-avg:faastroute && docker push $ID/wage-avg:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-format -t $ID/wage-format:faastroute && docker push $ID/wage-format:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-merit -t $ID/wage-merit:faastroute && docker push $ID/wage-merit:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-stats -t $ID/wage-stats:faastroute && docker push $ID/wage-stats:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-sum -t $ID/wage-sum:faastroute && docker push $ID/wage-sum:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-validator -t $ID/wage-validator:faastroute && docker push $ID/wage-validator:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-write-merit -t $ID/wage-write-merit:faastroute && docker push $ID/wage-write-merit:faastroute
docker build ./benchmarks/faastroute/original/wage-pay/wage-write-raw -t $ID/wage-write-raw:faastroute && docker push $ID/wage-write-raw:faastroute
docker image prune -a -f
