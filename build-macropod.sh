#!/bin/bash
ID=${1:-arearnest7}
docker build ./macropod-base/macropod-go -t $ID/macropod-go:latest && docker push $ID/macropod-go:latest
docker build ./macropod-base/macropod-node -t $ID/macropod-node:latest && docker push $ID/macropod-node:latest
docker build ./macropod-base/macropod-python -t $ID/macropod-python:latest && docker push $ID/macropod-python:latest
docker image prune -a -f
docker build ./benchmarks/micro/micro-rpc-a-macropod -t $ID/micro-rpc-a:macropod && docker push $ID/micro-rpc-a:macropod
docker build ./benchmarks/micro/micro-rpc-b-macropod -t $ID/micro-rpc-b:macropod && docker push $ID/micro-rpc-b:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/feature-generation/feature-extractor -t $ID/feature-extractor:macropod && docker push $ID/feature-extractor:macropod
docker build ./benchmarks/macropod/original/feature-generation/feature-orchestrator -t $ID/feature-orchestrator:macropod && docker push $ID/feature-orchestrator:macropod
docker build ./benchmarks/macropod/original/feature-generation/feature-reducer -t $ID/feature-reducer:macropod && docker push $ID/feature-reducer:macropod
docker build ./benchmarks/macropod/original/feature-generation/feature-status -t $ID/feature-status:macropod && docker push $ID/feature-status:macropod
docker build ./benchmarks/macropod/original/feature-generation/feature-wait -t $ID/feature-wait:macropod && docker push $ID/feature-wait:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:macropod && docker push $ID/hotel-frontend:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-geo -t $ID/hotel-geo:macropod && docker push $ID/hotel-geo:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-profile -t $ID/hotel-profile:macropod && docker push $ID/hotel-profile:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-rate -t $ID/hotel-rate:macropod && docker push $ID/hotel-rate:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:macropod && docker push $ID/hotel-recommend:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:macropod && docker push $ID/hotel-reserve:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-search -t $ID/hotel-search:macropod && docker push $ID/hotel-search:macropod
docker build ./benchmarks/macropod/original/hotel-app/hotel-user -t $ID/hotel-user:macropod && docker push $ID/hotel-user:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/pipelined-processing/pipelined-checksum -t $ID/pipelined-checksum:macropod && docker push $ID/pipelined-checksum:macropod
docker build ./benchmarks/macropod/original/pipelined-processing/pipelined-encrypt -t $ID/pipelined-encrypt:macropod && docker push $ID/pipelined-encrypt:macropod
docker build ./benchmarks/macropod/original/pipelined-processing/pipelined-main -t $ID/pipelined-main:macropod && docker push $ID/pipelined-main:macropod
docker build ./benchmarks/macropod/original/pipelined-processing/pipelined-zip -t $ID/pipelined-zip:macropod && docker push $ID/pipelined-zip:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-cfail -t $ID/sentiment-cfail:macropod && docker push $ID/sentiment-cfail:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-main -t $ID/sentiment-main:macropod && docker push $ID/sentiment-main:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-product-result -t $ID/sentiment-product-result:macropod && docker push $ID/sentiment-product-result:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-read-csv -t $ID/sentiment-read-csv:macropod && docker push $ID/sentiment-read-csv:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-service-sentiment -t $ID/sentiment-service-sentiment:macropod && docker push $ID/sentiment-service-sentiment:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-sns -t $ID/sentiment-sns:macropod && docker push $ID/sentiment-sns:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-db -t $ID/sentiment-db:macropod && docker push $ID/sentiment-db:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-product-or-service -t $ID/sentiment-product-or-service:macropod && docker push $ID/sentiment-product-or-service:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-product-sentiment -t $ID/sentiment-product-sentiment:macropod && docker push $ID/sentiment-product-sentiment:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-service-result -t $ID/sentiment-service-result:macropod && docker push $ID/sentiment-service-result:macropod
docker build ./benchmarks/macropod/original/sentiment-analysis/sentiment-sfail -t $ID/sentiment-sfail:macropod && docker push $ID/sentiment-sfail:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/serverless-election/election-gateway -t $ID/election-gateway:macropod && docker push $ID/election-gateway:macropod
docker build ./benchmarks/macropod/original/serverless-election/election-get-results -t $ID/election-get-results:macropod && docker push $ID/election-get-results:macropod
docker build ./benchmarks/macropod/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:macropod && docker push $ID/election-vote-enqueuer:macropod
docker build ./benchmarks/macropod/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:macropod && docker push $ID/election-vote-processor:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/video-analytics/video-decoder -t $ID/video-decoder:macropod && docker push $ID/video-decoder:macropod
docker build ./benchmarks/macropod/original/video-analytics/video-recog -t $ID/video-recog:macropod && docker push $ID/video-recog:macropod
docker build ./benchmarks/macropod/original/video-analytics/video-streaming -t $ID/video-streaming:macropod && docker push $ID/video-streaming:macropod
docker image prune -a -f
docker build ./benchmarks/macropod/original/wage-pay/wage-avg -t $ID/wage-avg:macropod && docker push $ID/wage-avg:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-format -t $ID/wage-format:macropod && docker push $ID/wage-format:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-merit -t $ID/wage-merit:macropod && docker push $ID/wage-merit:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-stats -t $ID/wage-stats:macropod && docker push $ID/wage-stats:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-sum -t $ID/wage-sum:macropod && docker push $ID/wage-sum:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-validator -t $ID/wage-validator:macropod && docker push $ID/wage-validator:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-write-merit -t $ID/wage-write-merit:macropod && docker push $ID/wage-write-merit:macropod
docker build ./benchmarks/macropod/original/wage-pay/wage-write-raw -t $ID/wage-write-raw:macropod && docker push $ID/wage-write-raw:macropod
docker image prune -a -f