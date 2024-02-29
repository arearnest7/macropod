#!/bin/bash
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/feature-generation/feature-extractor -t arearnest7/feature-extractor:faastroute && docker push arearnest7/feature-extractor:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/feature-generation/feature-orchestrator -t arearnest7/feature-orchestrator:faastroute && docker push arearnest7/feature-orchestrator:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/feature-generation/feature-reducer -t arearnest7/feature-reducer:faastroute && docker push arearnest7/feature-reducer:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/feature-generation/feature-status -t arearnest7/feature-status:faastroute && docker push arearnest7/feature-status:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/feature-generation/feature-wait -t arearnest7/feature-wait:faastroute && docker push arearnest7/feature-wait:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-frontend -t arearnest7/hotel-frontend:faastroute && docker push arearnest7/hotel-frontend:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-geo -t arearnest7/hotel-geo:faastroute && docker push arearnest7/hotel-geo:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-profile -t arearnest7/hotel-profile:faastroute && docker push arearnest7/hotel-profile:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-rate -t arearnest7/hotel-rate:faastroute && docker push arearnest7/hotel-rate:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-recommend -t arearnest7/hotel-recommend:faastroute && docker push arearnest7/hotel-recommend:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-reserve -t arearnest7/hotel-reserve:faastroute && docker push arearnest7/hotel-reserve:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-search -t arearnest7/hotel-search:faastroute && docker push arearnest7/hotel-search:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/hotel-app/hotel-user -t arearnest7/hotel-user:faastroute && docker push arearnest7/hotel-user:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/pipelined-processing/pipelined-checksum -t arearnest7/pipelined-checksum:faastroute && docker push arearnest7/pipelined-checksum:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/pipelined-processing/pipelined-encrypt -t arearnest7/pipelined-encrypt:faastroute && docker push arearnest7/pipelined-encrypt:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/pipelined-processing/pipelined-main -t arearnest7/pipelined-main:faastroute && docker push arearnest7/pipelined-main:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/pipelined-processing/pipelined-zip -t arearnest7/pipelined-zip:faastroute && docker push arearnest7/pipelined-zip:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-cfail -t arearnest7/sentiment-cfail:faastroute && docker push arearnest7/sentiment-cfail:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-main -t arearnest7/sentiment-main:faastroute && docker push arearnest7/sentiment-main:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-product-result -t arearnest7/sentiment-product-result:faastroute && docker push arearnest7/sentiment-product-result:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-read-csv -t arearnest7/sentiment-read-csv:faastroute && docker push arearnest7/sentiment-read-csv:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-service-sentiment -t arearnest7/sentiment-service-sentiment:faastroute && docker push arearnest7/sentiment-service-sentiment:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-sns -t arearnest7/sentiment-sns:faastroute && docker push arearnest7/sentiment-sns:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-db -t arearnest7/sentiment-db:faastroute && docker push arearnest7/sentiment-db:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-product-or-service -t arearnest7/sentiment-product-or-service:faastroute && docker push arearnest7/sentiment-product-or-service:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-product-sentiment -t arearnest7/sentiment-product-sentiment:faastroute && docker push arearnest7/sentiment-product-sentiment:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-service-result -t arearnest7/sentiment-service-result:faastroute && docker push arearnest7/sentiment-service-result:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/sentiment-analysis/sentiment-sfail -t arearnest7/sentiment-sfail:faastroute && docker push arearnest7/sentiment-sfail:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/serverless-election/election-gateway -t arearnest7/election-gateway:faastroute && docker push arearnest7/election-gateway:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/serverless-election/election-get-results -t arearnest7/election-get-results:faastroute && docker push arearnest7/election-get-results:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/serverless-election/election-vote-enqueuer -t arearnest7/election-vote-enqueuer:faastroute && docker push arearnest7/election-vote-enqueuer:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/serverless-election/election-vote-processor -t arearnest7/election-vote-processor:faastroute && docker push arearnest7/election-vote-processor:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/video-analytics/video-decoder -t arearnest7/video-decoder:faastroute && docker push arearnest7/video-decoder:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/video-analytics/video-recog -t arearnest7/video-recog:faastroute && docker push arearnest7/video-recog:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/video-analytics/video-streaming -t arearnest7/video-streaming:faastroute && docker push arearnest7/video-streaming:faastroute
docker image prune -a -f
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-avg -t arearnest7/wage-avg:faastroute && docker push arearnest7/wage-avg:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-format -t arearnest7/wage-format:faastroute && docker push arearnest7/wage-format:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-merit -t arearnest7/wage-merit:faastroute && docker push arearnest7/wage-merit:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-stats -t arearnest7/wage-stats:faastroute && docker push arearnest7/wage-stats:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-sum -t arearnest7/wage-sum:faastroute && docker push arearnest7/wage-sum:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-validator -t arearnest7/wage-validator:faastroute && docker push arearnest7/wage-validator:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-write-merit -t arearnest7/wage-write-merit:faastroute && docker push arearnest7/wage-write-merit:faastroute
docker build ./serverless-workflow-topology-survey/benchmarks/faastroute/original/wage-pay/wage-write-raw -t arearnest7/wage-write-raw:faastroute && docker push arearnest7/wage-write-raw:faastroute
docker image prune -a -f
