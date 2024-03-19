#!/bin/bash
ID=${1:-arearnest7}
docker build ./benchmarks/oci/original/feature-generation/feature-extractor -t $ID/feature-extractor:oci && docker push $ID/feature-extractor:oci
docker build ./benchmarks/oci/original/feature-generation/feature-orchestrator -t $ID/feature-orchestrator:oci && docker push $ID/feature-orchestrator:oci
docker build ./benchmarks/oci/original/feature-generation/feature-reducer -t $ID/feature-reducer:oci && docker push $ID/feature-reducer:oci
docker build ./benchmarks/oci/original/feature-generation/feature-status -t $ID/feature-status:oci && docker push $ID/feature-status:oci
docker build ./benchmarks/oci/original/feature-generation/feature-wait -t $ID/feature-wait:oci && docker push $ID/feature-wait:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:oci && docker push $ID/hotel-frontend:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-geo -t $ID/hotel-geo:oci && docker push $ID/hotel-geo:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-profile -t $ID/hotel-profile:oci && docker push $ID/hotel-profile:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-rate -t $ID/hotel-rate:oci && docker push $ID/hotel-rate:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:oci && docker push $ID/hotel-recommend:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:oci && docker push $ID/hotel-reserve:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-search -t $ID/hotel-search:oci && docker push $ID/hotel-search:oci
docker build ./benchmarks/oci/original/hotel-app/hotel-user -t $ID/hotel-user:oci && docker push $ID/hotel-user:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/pipelined-processing/pipelined-checksum -t $ID/pipelined-checksum:oci && docker push $ID/pipelined-checksum:oci
docker build ./benchmarks/oci/original/pipelined-processing/pipelined-encrypt -t $ID/pipelined-encrypt:oci && docker push $ID/pipelined-encrypt:oci
docker build ./benchmarks/oci/original/pipelined-processing/pipelined-main -t $ID/pipelined-main:oci && docker push $ID/pipelined-main:oci
docker build ./benchmarks/oci/original/pipelined-processing/pipelined-zip -t $ID/pipelined-zip:oci && docker push $ID/pipelined-zip:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-cfail -t $ID/sentiment-cfail:oci && docker push $ID/sentiment-cfail:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-main -t $ID/sentiment-main:oci && docker push $ID/sentiment-main:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-product-result -t $ID/sentiment-product-result:oci && docker push $ID/sentiment-product-result:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-read-csv -t $ID/sentiment-read-csv:oci && docker push $ID/sentiment-read-csv:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-service-sentiment -t $ID/sentiment-service-sentiment:oci && docker push $ID/sentiment-service-sentiment:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-sns -t $ID/sentiment-sns:oci && docker push $ID/sentiment-sns:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-db -t $ID/sentiment-db:oci && docker push $ID/sentiment-db:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-product-or-service -t $ID/sentiment-product-or-service:oci && docker push $ID/sentiment-product-or-service:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-product-sentiment -t $ID/sentiment-product-sentiment:oci && docker push $ID/sentiment-product-sentiment:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-service-result -t $ID/sentiment-service-result:oci && docker push $ID/sentiment-service-result:oci
docker build ./benchmarks/oci/original/sentiment-analysis/sentiment-sfail -t $ID/sentiment-sfail:oci && docker push $ID/sentiment-sfail:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/serverless-election/election-gateway -t $ID/election-gateway:oci && docker push $ID/election-gateway:oci
docker build ./benchmarks/oci/original/serverless-election/election-get-results -t $ID/election-get-results:oci && docker push $ID/election-get-results:oci
docker build ./benchmarks/oci/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:oci && docker push $ID/election-vote-enqueuer:oci
docker build ./benchmarks/oci/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:oci && docker push $ID/election-vote-processor:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/video-analytics/video-decoder -t $ID/video-decoder:oci && docker push $ID/video-decoder:oci
docker build ./benchmarks/oci/original/video-analytics/video-recog -t $ID/video-recog:oci && docker push $ID/video-recog:oci
docker build ./benchmarks/oci/original/video-analytics/video-streaming -t $ID/video-streaming:oci && docker push $ID/video-streaming:oci
docker image prune -a -f
docker build ./benchmarks/oci/original/wage-pay/wage-avg -t $ID/wage-avg:oci && docker push $ID/wage-avg:oci
docker build ./benchmarks/oci/original/wage-pay/wage-format -t $ID/wage-format:oci && docker push $ID/wage-format:oci
docker build ./benchmarks/oci/original/wage-pay/wage-merit -t $ID/wage-merit:oci && docker push $ID/wage-merit:oci
docker build ./benchmarks/oci/original/wage-pay/wage-stats -t $ID/wage-stats:oci && docker push $ID/wage-stats:oci
docker build ./benchmarks/oci/original/wage-pay/wage-sum -t $ID/wage-sum:oci && docker push $ID/wage-sum:oci
docker build ./benchmarks/oci/original/wage-pay/wage-validator -t $ID/wage-validator:oci && docker push $ID/wage-validator:oci
docker build ./benchmarks/oci/original/wage-pay/wage-write-merit -t $ID/wage-write-merit:oci && docker push $ID/wage-write-merit:oci
docker build ./benchmarks/oci/original/wage-pay/wage-write-raw -t $ID/wage-write-raw:oci && docker push $ID/wage-write-raw:oci
docker image prune -a -f
