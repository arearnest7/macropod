#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-macropod}
DOCKER_BUILDKIT=1 docker build base/macropod-ingress -t $ID/macropod-ingress:$TAG && docker push $ID/macropod-ingress:$TAG
DOCKER_BUILDKIT=1 docker build base/macropod-deployer -t $ID/macropod-deployer:$TAG && docker push $ID/macropod-deployer:$TAG

docker build base/macropod-go -t $ID/macropod-go:$TAG && docker push $ID/macropod-go:$TAG
docker build base/macropod-node -t $ID/macropod-node:$TAG && docker push $ID/macropod-node:$TAG
docker build base/macropod-python -t $ID/macropod-python:$TAG && docker push $ID/macropod-python:$TAG

docker build benchmarks/macropod/unified/hotel-app/hotel-unified -t $ID/hotel-unified:$TAG && docker push $ID/hotel-unified:$TAG
docker build benchmarks/macropod/unified/pipelined-processing/pipelined-unified -t $ID/pipelined-unified:$TAG && docker push $ID/pipelined-unified:$TAG
docker build benchmarks/macropod/unified/sentiment-analysis/sentiment-unified -t $ID/sentiment-unified:$TAG && docker push $ID/sentiment-unified:$TAG
docker build benchmarks/macropod/unified/serverless-election/election-unified -t $ID/election-unified:$TAG && docker push $ID/election-unified:$TAG
docker build benchmarks/macropod/unified/video-analytics/video-unified -t $ID/video-unified:$TAG && docker push $ID/video-unified:$TAG
docker build benchmarks/macropod/unified/wage-pay/wage-unified -t $ID/wage-unified:$TAG && docker push $ID/wage-unified:$TAG

docker build benchmarks/macropod/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:$TAG && docker push $ID/hotel-frontend:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-geo -t $ID/hotel-geo:$TAG && docker push $ID/hotel-geo:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-profile -t $ID/hotel-profile:$TAG && docker push $ID/hotel-profile:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-rate -t $ID/hotel-rate:$TAG && docker push $ID/hotel-rate:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:$TAG && docker push $ID/hotel-recommend:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:$TAG && docker push $ID/hotel-reserve:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-search -t $ID/hotel-search:$TAG && docker push $ID/hotel-search:$TAG
docker build benchmarks/macropod/original/hotel-app/hotel-user -t $ID/hotel-user:$TAG && docker push $ID/hotel-user:$TAG
docker build benchmarks/macropod/original/pipelined-processing/pipelined-checksum -t $ID/pipelined-checksum:$TAG && docker push $ID/pipelined-checksum:$TAG
docker build benchmarks/macropod/original/pipelined-processing/pipelined-encrypt -t $ID/pipelined-encrypt:$TAG && docker push $ID/pipelined-encrypt:$TAG
docker build benchmarks/macropod/original/pipelined-processing/pipelined-main -t $ID/pipelined-main:$TAG && docker push $ID/pipelined-main:$TAG
docker build benchmarks/macropod/original/pipelined-processing/pipelined-zip -t $ID/pipelined-zip:$TAG && docker push $ID/pipelined-zip:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-cfail -t $ID/sentiment-cfail:$TAG && docker push $ID/sentiment-cfail:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-main -t $ID/sentiment-main:$TAG && docker push $ID/sentiment-main:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-product-result -t $ID/sentiment-product-result:$TAG && docker push $ID/sentiment-product-result:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-read-csv -t $ID/sentiment-read-csv:$TAG && docker push $ID/sentiment-read-csv:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-service-sentiment -t $ID/sentiment-service-sentiment:$TAG && docker push $ID/sentiment-service-sentiment:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-sns -t $ID/sentiment-sns:$TAG && docker push $ID/sentiment-sns:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-db -t $ID/sentiment-db:$TAG && docker push $ID/sentiment-db:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-product-or-service -t $ID/sentiment-product-or-service:$TAG && docker push $ID/sentiment-product-or-service:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-product-sentiment -t $ID/sentiment-product-sentiment:$TAG && docker push $ID/sentiment-product-sentiment:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-service-result -t $ID/sentiment-service-result:$TAG && docker push $ID/sentiment-service-result:$TAG
docker build benchmarks/macropod/original/sentiment-analysis/sentiment-sfail -t $ID/sentiment-sfail:$TAG && docker push $ID/sentiment-sfail:$TAG
docker build benchmarks/macropod/original/serverless-election/election-gateway -t $ID/election-gateway:$TAG && docker push $ID/election-gateway:$TAG
docker build benchmarks/macropod/original/serverless-election/election-get-results -t $ID/election-get-results:$TAG && docker push $ID/election-get-results:$TAG
docker build benchmarks/macropod/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:$TAG && docker push $ID/election-vote-enqueuer:$TAG
docker build benchmarks/macropod/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:$TAG && docker push $ID/election-vote-processor:$TAG
docker build benchmarks/macropod/original/video-analytics/video-decoder -t $ID/video-decoder:$TAG && docker push $ID/video-decoder:$TAG
docker build benchmarks/macropod/original/video-analytics/video-recog -t $ID/video-recog:$TAG && docker push $ID/video-recog:$TAG
docker build benchmarks/macropod/original/video-analytics/video-streaming -t $ID/video-streaming:$TAG && docker push $ID/video-streaming:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-avg -t $ID/wage-avg:$TAG && docker push $ID/wage-avg:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-format -t $ID/wage-format:$TAG && docker push $ID/wage-format:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-merit -t $ID/wage-merit:$TAG && docker push $ID/wage-merit:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-stats -t $ID/wage-stats:$TAG && docker push $ID/wage-stats:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-sum -t $ID/wage-sum:$TAG && docker push $ID/wage-sum:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-validator -t $ID/wage-validator:$TAG && docker push $ID/wage-validator:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-write-merit -t $ID/wage-write-merit:$TAG && docker push $ID/wage-write-merit:$TAG
docker build benchmarks/macropod/original/wage-pay/wage-write-raw -t $ID/wage-write-raw:$TAG && docker push $ID/wage-write-raw:$TAG
