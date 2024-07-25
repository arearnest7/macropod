#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
./collect-kn.sh $HOST full-election election-full 1 ../../payloads/election.json 500
./collect-kn.sh $HOST full-feature feature-full 1 ../../payloads/feature.json 500
./collect-kn.sh $HOST full-hotel hotel-full 1 ../../payloads/hotel.json 500
./collect-kn.sh $HOST full-pipelined pipelined-full 1 ../../payloads/pipelined.json 500
./collect-kn.sh $HOST full-sentiment sentiment-full 1 ../../payloads/sentiment.json 500
./collect-kn.sh $HOST full-video video-full 1 ../../payloads/video.json 500
./collect-kn.sh $HOST full-wage wage-full 1 ../../payloads/wage.json 500
./collect-kn.sh $HOST original-election election-gateway 1 ../../payloads/election.json 500
./collect-kn.sh $HOST original-feature feature-orchestrator 1 ../../payloads/feature.json 500
./collect-kn.sh $HOST original-hotel hotel-frontend 1 ../../payloads/hotel.json 500
./collect-kn.sh $HOST original-pipelined pipelined-main 1 ../../payloads/pipelined.json 500
./collect-kn.sh $HOST original-sentiment sentiment-main 1 ../../payloads/sentiment.json 500
./collect-kn.sh $HOST original-video video-streaming 1 ../../payloads/video.json 500
./collect-kn.sh $HOST original-wage wage-validator 1 ../../payloads/wage.json 500
./collect-kn.sh $HOST partial-election election-gateway-vevp 1 ../../payloads/election.json 500
./collect-kn.sh $HOST partial-feature feature-orchestrator-wsr 1 ../../payloads/feature.json 500
./collect-kn.sh $HOST partial-hotel hotel-frontend-spgr 1 ../../payloads/hotel.json 500
./collect-kn.sh $HOST partial-pipelined pipelined-main-partial 1 ../../payloads/pipelined.json 500
./collect-kn.sh $HOST partial-sentiment sentiment-main-rcposc 1 ../../payloads/sentiment.json 500
./collect-kn.sh $HOST partial-video video-streaming-d 1 ../../payloads/video.json 500
./collect-kn.sh $HOST partial-wage wage-validator-fw 1 ../../payloads/wage.json 500
./collect-yaml.sh election 1 ../payloads/election.json 500 "election-gateway election-get-results election-vote-enqueuer election-vote-processor"
./collect-yaml.sh feature 1 ../payloads/feature.json 500 "feature-orchestrator feature-extractor feature-wait feature-status feature-reducer"
./collect-yaml.sh hotel 1 ../payloads/hotel.json 500 "hotel-frontend hotel-search hotel-recommend hotel-reserve hotel-user hotel-geo hotel-profile hotel-rate"
./collect-yaml.sh pipelined 1 ../payloads/pipelined.json 500 "pipelined-main pipelined-main-2 pipelined-main-3 pipelined-checksum pipelined-zip pipelined-encrypt"
./collect-yaml.sh sentiment 1 ../payloads/sentiment.json 500 "sentiment-cfail sentiment-db sentiment-main sentiment-product-or-service sentiment-product-result sentiment-product-sentiment sentiment-read-csv sentiment-service-result sentiment-service-sentiment sentiment-sfail sentiment-sns"
./collect-yaml.sh video 1 ../payloads/video.json 500 "video-streaming video-decoder video-recog"
./collect-yaml.sh wage 1 ../payloads/wage.json 500 "wage-validator wage-format wage-write-raw wage-stats wage-sum wage-avg wage-merit wage-write-merit"
mkdir $DIR
mv *.csv $DIR
