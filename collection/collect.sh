#!/bin/bash
HOST=${1:-127.0.0.1}
REDIS=${2:-127.0.0.1}
PASSWORD=${3:-password}
C=${4:-1}
N=${5:-10000}
DIR=${6:-results}
./collect-invoke.sh $HOST full-election wob election-full $REDIS $PASSWORD $C ../payloads/election.json $N
./collect-invoke.sh $HOST full-feature wob feature-full $REDIS $PASSWORD $C ../payloads/feature.json $N
./collect-invoke.sh $HOST full-hotel wob hotel-full $REDIS $PASSWORD $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST full-pipelined wob pipelined-full $REDIS $PASSWORD $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST full-sentiment wob sentiment-full $REDIS $PASSWORD $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST full-video wob video-full $REDIS $PASSWORD $C ../payloads/video.json $N
./collect-invoke.sh $HOST full-wage wob wage-full $REDIS $PASSWORD $C ../payloads/wage.json $N
./collect-invoke.sh $HOST original-election wob election-gateway $REDIS $PASSWORD $C ../payloads/election.json $N
./collect-invoke.sh $HOST original-feature wob feature-orchestrator $REDIS $PASSWORD $C ../payloads/feature.json $N
./collect-invoke.sh $HOST original-hotel wob hotel-frontend $REDIS $PASSWORD $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST original-pipelined wob pipelined-main $REDIS $PASSWORD $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST original-sentiment wob sentiment-main $REDIS $PASSWORD $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST original-video wob video-streaming $REDIS $PASSWORD $C ../payloads/video.json $N
./collect-invoke.sh $HOST original-wage wob wage-validator $REDIS $PASSWORD $C ../payloads/wage.json $N
./collect-invoke.sh $HOST partial-election wob election-gateway $REDIS $PASSWORD $C ../payloads/election.json $N
./collect-invoke.sh $HOST partial-feature wob feature-orchestrator-wsr $REDIS $PASSWORD $C ../payloads/feature.json $N
./collect-invoke.sh $HOST partial-hotel wob hotel-frontend-spgr $REDIS $PASSWORD $C ../payloads/hotel.json $N
./collect-invoke.sh $HOST partial-pipelined wob pipelined-main-partial $REDIS $PASSWORD $C ../payloads/pipelined.json $N
./collect-invoke.sh $HOST partial-sentiment wob sentiment-main-rcposc $REDIS $PASSWORD $C ../payloads/sentiment.json $N
./collect-invoke.sh $HOST partial-video wob video-streaming-d $REDIS $PASSWORD $C ../payloads/video.json $N
./collect-invoke.sh $HOST partial-wage wob wage-validator-fw $REDIS $PASSWORD $C ../payloads/wage.json $N
./collect-yaml.sh election $REDIS $PASSWORD $C ../payloads/election.json $N
./collect-yaml.sh feature $REDIS $PASSWORD $C ../payloads/feature.json $N
./collect-yaml.sh hotel $REDIS $PASSWORD $C ../payloads/hotel.json $N
./collect-yaml.sh pipelined $REDIS $PASSWORD $C ../payloads/pipelined.json $N
./collect-yaml.sh sentiment $REDIS $PASSWORD $C ../payloads/sentiment.json $N
./collect-yaml.sh video $REDIS $PASSWORD $C ../payloads/video.json $N
./collect-yaml.sh wage $REDIS $PASSWORD $C ../payloads/wage.json $N
mkdir $DIR
mv *.csv $DIR
