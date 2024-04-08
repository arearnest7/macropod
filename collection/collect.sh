#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
DIR=${4:-results}
./collect-invoke.sh $HOST full-election wob election-full $C ../payloads/election.json $N
./collect-invoke.sh $HOST full-feature wob feature-full $C ../payloads/feature.json $N
./collect-invoke.sh $HOST full-hotel wob hotel-full $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST full-pipelined wob pipelined-full $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST full-sentiment wob sentiment-full $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST full-video wob video-full $C ../payloads/video.json $N
./collect-invoke.sh $HOST full-wage wob wage-full $C ../payloads/wage.json $N
./collect-invoke.sh $HOST original-election wob election-gateway $C ../payloads/election.json $N
./collect-invoke.sh $HOST original-feature wob feature-orchestrator $C ../payloads/feature.json $N
./collect-invoke.sh $HOST original-hotel wob hotel-frontend $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST original-pipelined wob pipelined-main $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST original-sentiment wob sentiment-main $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST original-video wob video-streaming $C ../payloads/video.json $N
./collect-invoke.sh $HOST original-wage wob wage-validator $C ../payloads/wage.json $N
./collect-invoke.sh $HOST partial-election wob election-gateway $C ../payloads/election.json $N
./collect-invoke.sh $HOST partial-feature wob feature-orchestrator-wsr $C ../payloads/feature.json $N
./collect-invoke.sh $HOST partial-hotel wob hotel-frontend-spgr $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST partial-pipelined wob pipelined-main-partial $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST partial-sentiment wob sentiment-main-rcposc $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST partial-video wob video-streaming-d $C ../payloads/video.json $N
./collect-invoke.sh $HOST partial-wage wob wage-validator-fw $C ../payloads/wage.json $N
./collect-yaml.sh election $C ../payloads/election.json $N
./collect-yaml.sh feature $C ../payloads/feature.json $N
./collect-yaml.sh hotel $C ../payloads/hotel.json $N
./collect-yaml.sh pipelined $C ../payloads/pipelined.json $N
./collect-yaml.sh sentiment $C ../payloads/sentiment.json $N
./collect-yaml.sh video $C ../payloads/video.json $N
./collect-yaml.sh wage $C ../payloads/wage.json $N
mkdir $DIR
mv *.csv $DIR
