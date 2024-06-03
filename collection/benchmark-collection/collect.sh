#!/bin/bash
HOST=${1:-127.0.0.1}
C=${2:-1}
N=${3:-10000}
DIR=${4:-results}
./collect-kn.sh $HOST full-election wob election-full $C ../payloads/election.json $N
./collect-kn.sh $HOST full-feature wob feature-full $C ../payloads/feature.json $N
./collect-kn.sh $HOST full-hotel wob hotel-full $C ../payloads/hotel.json $N
./collect-kn.sh $HOST full-pipelined wob pipelined-full $C ../payloads/pipelined.json $N
./collect-kn.sh $HOST full-sentiment wob sentiment-full $C ../payloads/sentiment.json $N
./collect-kn.sh $HOST full-video wob video-full $C ../payloads/video.json $N
./collect-kn.sh $HOST full-wage wob wage-full $C ../payloads/wage.json $N
./collect-kn.sh $HOST original-election wob election-gateway $C ../payloads/election.json $N
./collect-kn.sh $HOST original-feature wob feature-orchestrator $C ../payloads/feature.json $N
./collect-kn.sh $HOST original-hotel wob hotel-frontend $C ../payloads/hotel.json $N
./collect-kn.sh $HOST original-pipelined wob pipelined-main $C ../payloads/pipelined.json $N
./collect-kn.sh $HOST original-sentiment wob sentiment-main $C ../payloads/sentiment.json $N
./collect-kn.sh $HOST original-video wob video-streaming $C ../payloads/video.json $N
./collect-kn.sh $HOST original-wage wob wage-validator $C ../payloads/wage.json $N
./collect-kn.sh $HOST partial-election wob election-gateway-vevp $C ../payloads/election.json $N
./collect-kn.sh $HOST partial-feature wob feature-orchestrator-wsr $C ../payloads/feature.json $N
./collect-kn.sh $HOST partial-hotel wob hotel-frontend-spgr $C ../payloads/hotel.json $N
./collect-kn.sh $HOST partial-pipelined wob pipelined-main-partial $C ../payloads/pipelined.json $N
./collect-kn.sh $HOST partial-sentiment wob sentiment-main-rcposc $C ../payloads/sentiment.json $N
./collect-kn.sh $HOST partial-video wob video-streaming-d $C ../payloads/video.json $N
./collect-kn.sh $HOST partial-wage wob wage-validator-fw $C ../payloads/wage.json $N
./collect-macropod.sh election $C ../payloads/election.json $N
./collect-macropod.sh feature $C ../payloads/feature.json $N
./collect-macropod.sh hotel $C ../payloads/hotel.json $N
./collect-macropod.sh pipelined $C ../payloads/pipelined.json $N
./collect-macropod.sh sentiment $C ../payloads/sentiment.json $N
./collect-macropod.sh video $C ../payloads/video.json $N
./collect-macropod.sh wage $C ../payloads/wage.json $N
mkdir $DIR
mv *.csv $DIR
