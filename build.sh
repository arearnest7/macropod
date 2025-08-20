#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
sudo docker buildx build base/macropod-ingress -t $ID/macropod-ingress:$TAG && sudo docker push $ID/macropod-ingress:$TAG
sudo docker buildx build base/macropod-deployer -t $ID/macropod-deployer:$TAG && sudo docker push $ID/macropod-deployer:$TAG
sudo docker buildx build base/macropod-eval -t $ID/macropod-eval:$TAG && sudo docker push $ID/macropod-eval:$TAG
sudo docker buildx build base/macropod-logger -t $ID/macropod-logger:$TAG && sudo docker push $ID/macropod-logger:$TAG
sudo docker buildx build base/macropod-metrics -t $ID/macropod-metrics:$TAG && sudo docker push $ID/macropod-metrics:$TAG

sudo docker buildx build base/macropod-go -t $ID/macropod-go:$TAG && sudo docker push $ID/macropod-go:$TAG
sudo docker buildx build base/macropod-node -t $ID/macropod-node:$TAG && sudo docker push $ID/macropod-node:$TAG
sudo docker buildx build base/macropod-python -t $ID/macropod-python:$TAG && sudo docker push $ID/macropod-python:$TAG

sudo docker buildx build benchmarks/macropod/unified/hotel-app/hotel-unified -t $ID/hotel-unified:$TAG && sudo docker push $ID/hotel-unified:$TAG
sudo docker buildx build benchmarks/macropod/unified/pipelined-processing/pipelined-unified -t $ID/pipelined-unified:$TAG && sudo docker push $ID/pipelined-unified:$TAG
sudo docker buildx build benchmarks/macropod/unified/sentiment-analysis/sentiment-unified -t $ID/sentiment-unified:$TAG && sudo docker push $ID/sentiment-unified:$TAG
sudo docker buildx build benchmarks/macropod/unified/serverless-election/election-unified -t $ID/election-unified:$TAG && sudo docker push $ID/election-unified:$TAG
sudo docker buildx build benchmarks/macropod/unified/video-analytics/video-unified -t $ID/video-unified:$TAG && sudo docker push $ID/video-unified:$TAG
sudo docker buildx build benchmarks/macropod/unified/wage-pay/wage-unified -t $ID/wage-unified:$TAG && sudo docker push $ID/wage-unified:$TAG

sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:$TAG && sudo docker push $ID/hotel-frontend:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-geo -t $ID/hotel-geo:$TAG && sudo docker push $ID/hotel-geo:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-profile -t $ID/hotel-profile:$TAG && sudo docker push $ID/hotel-profile:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-rate -t $ID/hotel-rate:$TAG && sudo docker push $ID/hotel-rate:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:$TAG && sudo docker push $ID/hotel-recommend:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:$TAG && sudo docker push $ID/hotel-reserve:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-search -t $ID/hotel-search:$TAG && sudo docker push $ID/hotel-search:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-user -t $ID/hotel-user:$TAG && sudo docker push $ID/hotel-user:$TAG
sudo docker buildx build benchmarks/macropod/original/pipelined-processing/pipelined-checksum -t $ID/pipelined-checksum:$TAG && sudo docker push $ID/pipelined-checksum:$TAG
sudo docker buildx build benchmarks/macropod/original/pipelined-processing/pipelined-encrypt -t $ID/pipelined-encrypt:$TAG && sudo docker push $ID/pipelined-encrypt:$TAG
sudo docker buildx build benchmarks/macropod/original/pipelined-processing/pipelined-main -t $ID/pipelined-main:$TAG && sudo docker push $ID/pipelined-main:$TAG
sudo docker buildx build benchmarks/macropod/original/pipelined-processing/pipelined-zip -t $ID/pipelined-zip:$TAG && sudo docker push $ID/pipelined-zip:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-cfail -t $ID/sentiment-cfail:$TAG && sudo docker push $ID/sentiment-cfail:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-main -t $ID/sentiment-main:$TAG && sudo docker push $ID/sentiment-main:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-product-result -t $ID/sentiment-product-result:$TAG && sudo docker push $ID/sentiment-product-result:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-read-csv -t $ID/sentiment-read-csv:$TAG && sudo docker push $ID/sentiment-read-csv:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-service-sentiment -t $ID/sentiment-service-sentiment:$TAG && sudo docker push $ID/sentiment-service-sentiment:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-sns -t $ID/sentiment-sns:$TAG && sudo docker push $ID/sentiment-sns:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-db -t $ID/sentiment-db:$TAG && sudo docker push $ID/sentiment-db:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-product-or-service -t $ID/sentiment-product-or-service:$TAG && sudo docker push $ID/sentiment-product-or-service:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-product-sentiment -t $ID/sentiment-product-sentiment:$TAG && sudo docker push $ID/sentiment-product-sentiment:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-service-result -t $ID/sentiment-service-result:$TAG && sudo docker push $ID/sentiment-service-result:$TAG
sudo docker buildx build benchmarks/macropod/original/sentiment-analysis/sentiment-sfail -t $ID/sentiment-sfail:$TAG && sudo docker push $ID/sentiment-sfail:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-gateway -t $ID/election-gateway:$TAG && sudo docker push $ID/election-gateway:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-get-results -t $ID/election-get-results:$TAG && sudo docker push $ID/election-get-results:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:$TAG && sudo docker push $ID/election-vote-enqueuer:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:$TAG && sudo docker push $ID/election-vote-processor:$TAG
sudo docker buildx build benchmarks/macropod/original/video-analytics/video-decoder -t $ID/video-decoder:$TAG && sudo docker push $ID/video-decoder:$TAG
sudo docker buildx build benchmarks/macropod/original/video-analytics/video-recog -t $ID/video-recog:$TAG && sudo docker push $ID/video-recog:$TAG
sudo docker buildx build benchmarks/macropod/original/video-analytics/video-streaming -t $ID/video-streaming:$TAG && sudo docker push $ID/video-streaming:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-avg -t $ID/wage-avg:$TAG && sudo docker push $ID/wage-avg:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-format -t $ID/wage-format:$TAG && sudo docker push $ID/wage-format:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-merit -t $ID/wage-merit:$TAG && sudo docker push $ID/wage-merit:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-stats -t $ID/wage-stats:$TAG && sudo docker push $ID/wage-stats:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-sum -t $ID/wage-sum:$TAG && sudo docker push $ID/wage-sum:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-validator -t $ID/wage-validator:$TAG && sudo docker push $ID/wage-validator:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-write-merit -t $ID/wage-write-merit:$TAG && sudo docker push $ID/wage-write-merit:$TAG
sudo docker buildx build benchmarks/macropod/original/wage-pay/wage-write-raw -t $ID/wage-write-raw:$TAG && sudo docker push $ID/wage-write-raw:$TAG
