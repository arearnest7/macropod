#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
DIR=${4:-results}
./collect-kn.sh $HOST full-election election-full $C ../../payloads/election.json $N
./collect-kn.sh $HOST full-feature feature-full $C ../../payloads/feature.json $N
./collect-kn.sh $HOST full-hotel hotel-full $C ../../payloads/hotel.json $N
./collect-kn.sh $HOST full-pipelined pipelined-full $C ../../payloads/pipelined.json $N
./collect-kn.sh $HOST full-sentiment sentiment-full $C ../../payloads/sentiment.json $N
./collect-kn.sh $HOST full-video video-full $C ../../payloads/video.json $N
./collect-kn.sh $HOST full-wage wage-full $C ../../payloads/wage.json $N
./collect-kn.sh $HOST original-election election-gateway $C ../../payloads/election.json $N
./collect-kn.sh $HOST original-feature feature-orchestrator $C ../../payloads/feature.json $N
./collect-kn.sh $HOST original-hotel hotel-frontend $C ../../payloads/hotel.json $N
./collect-kn.sh $HOST original-pipelined pipelined-main $C ../../payloads/pipelined.json $N
./collect-kn.sh $HOST original-sentiment sentiment-main $C ../../payloads/sentiment.json $N
./collect-kn.sh $HOST original-video video-streaming $C ../../payloads/video.json $N
./collect-kn.sh $HOST original-wage wage-validator $C ../../payloads/wage.json $N
./collect-kn.sh $HOST partial-election election-gateway-vevp $C ../../payloads/election.json $N
./collect-kn.sh $HOST partial-feature feature-orchestrator-wsr $C ../../payloads/feature.json $N
./collect-kn.sh $HOST partial-hotel hotel-frontend-spgr $C ../../payloads/hotel.json $N
./collect-kn.sh $HOST partial-pipelined pipelined-main-partial $C ../../payloads/pipelined.json $N
./collect-kn.sh $HOST partial-sentiment sentiment-main-rcposc $C ../../payloads/sentiment.json $N
./collect-kn.sh $HOST partial-video video-streaming-d $C ../../payloads/video.json $N
./collect-kn.sh $HOST partial-wage wage-validator-fw $C ../../payloads/wage.json $N
./collect-yaml.sh election $C ../payloads/election.json $N "election-gateway election-get-results election-vote-enqueuer election-vote-processor"
./collect-yaml.sh feature $C ../payloads/feature.json $N "feature-orchestrator feature-extractor feature-wait feature-status feature-reducer"
./collect-yaml.sh hotel $C ../payloads/hotel.json $N "hotel-frontend hotel-search hotel-recommend hotel-reserve hotel-user hotel-geo hotel-profile hotel-rate"
./collect-yaml.sh pipelined $C ../payloads/pipelined.json $N "pipelined-main pipelined-main-2 pipelined-main-3 pipelined-checksum pipelined-zip pipelined-encrypt"
./collect-yaml.sh sentiment $C ../payloads/sentiment.json $N "sentiment-cfail sentiment-db sentiment-main sentiment-product-or-service sentiment-product-result sentiment-product-sentiment sentiment-read-csv sentiment-service-result sentiment-service-sentiment sentiment-sfail sentiment-sns"
./collect-yaml.sh video $C ../payloads/video.json $N "video-streaming video-decoder video-recog"
./collect-yaml.sh wage $C ../payloads/wage.json $N "wage-validator wage-format wage-write-raw wage-stats wage-sum wage-avg wage-merit wage-write-merit"
mkdir $DIR
mv *.csv $DIR