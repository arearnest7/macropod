#!/bin/bash
cd ./tools/deploy-functions-kn/
kn func deploy --build=true --push=true --path ../../benchmarks/micro/micro-rpc-a
kn func deploy --build=true --push=true --path ../../benchmarks/micro/micro-rpc-b
kn func deploy --build=true --push=true --path ../../benchmarks/micro/micro-rpc-a-b
kn func delete micro-rpc-a
kn func delete micro-rpc-b
kn func delete micro-rpc-a-b
docker system prune -a -f
docker volume prune -f
./deploy-full-election.sh kn true true
../remove-functions-kn/remove-full-election.sh kn
./deploy-full-feature.sh kn true true
../remove-functions-kn/remove-full-feature.sh kn
./deploy-full-hotel.sh kn true true
../remove-functions-kn/remove-full-hotel.sh kn
./deploy-full-pipelined.sh kn true true
../remove-functions-kn/remove-full-pipelined.sh kn
./deploy-full-sentiment.sh kn true true
../remove-functions-kn/remove-full-sentiment.sh kn
./deploy-full-video.sh kn true true
../remove-functions-kn/remove-full-video.sh kn
./deploy-full-wage.sh kn true true
../remove-functions-kn/remove-full-wage.sh kn
docker system prune -a -f
docker volume prune -f
./deploy-original-election.sh kn true true
../remove-functions-kn/remove-original-election.sh kn
./deploy-original-feature.sh kn true true
../remove-functions-kn/remove-original-feature.sh kn
./deploy-original-hotel.sh kn true true
../remove-functions-kn/remove-original-hotel.sh kn
./deploy-original-pipelined.sh kn true true
../remove-functions-kn/remove-original-pipelined.sh kn
./deploy-original-sentiment.sh kn true true
../remove-functions-kn/remove-original-sentiment.sh kn
./deploy-original-video.sh kn true true
../remove-functions-kn/remove-original-video.sh kn
./deploy-original-wage.sh kn true true
../remove-functions-kn/remove-original-wage.sh kn
docker system prune -a -f
docker volume prune -f
./deploy-partial-election.sh kn true true
../remove-functions-kn/remove-partial-election.sh kn
./deploy-partial-feature.sh kn true true
../remove-functions-kn/remove-partial-feature.sh kn
./deploy-partial-hotel.sh kn true true
../remove-functions-kn/remove-partial-hotel.sh kn
./deploy-partial-pipelined.sh kn true true
../remove-functions-kn/remove-partial-pipelined.sh kn
./deploy-partial-sentiment.sh kn true true
../remove-functions-kn/remove-partial-sentiment.sh kn
./deploy-partial-video.sh kn true true
../remove-functions-kn/remove-partial-video.sh kn
./deploy-partial-wage.sh kn true true
../remove-functions-kn/remove-partial-wage.sh kn
docker system prune -a -f
docker volume prune -f
./deploy-full-election.sh wob true true
../remove-functions-kn/remove-full-election.sh wob
./deploy-full-feature.sh wob true true
../remove-functions-kn/remove-full-feature.sh wob
./deploy-full-hotel.sh wob true true
../remove-functions-kn/remove-full-hotel.sh wob
./deploy-full-pipelined.sh wob true true
../remove-functions-kn/remove-full-pipelined.sh wob
./deploy-full-sentiment.sh wob true true
../remove-functions-kn/remove-full-sentiment.sh wob
./deploy-full-video.sh wob true true
../remove-functions-kn/remove-full-video.sh wob
./deploy-full-wage.sh wob true true
../remove-functions-kn/remove-full-wage.sh wob
docker system prune -a -f
docker volume prune -f
./deploy-original-election.sh wob true true
../remove-functions-kn/remove-original-election.sh wob
./deploy-original-feature.sh wob true true
../remove-functions-kn/remove-original-feature.sh wob
./deploy-original-hotel.sh wob true true
../remove-functions-kn/remove-original-hotel.sh wob
./deploy-original-pipelined.sh wob true true
../remove-functions-kn/remove-original-pipelined.sh wob
./deploy-original-sentiment.sh wob true true
../remove-functions-kn/remove-original-sentiment.sh wob
./deploy-original-video.sh wob true true
../remove-functions-kn/remove-original-video.sh wob
./deploy-original-wage.sh wob true true
../remove-functions-kn/remove-original-wage.sh wob
docker system prune -a -f
docker volume prune -f
./deploy-partial-election.sh wob true true
../remove-functions-kn/remove-partial-election.sh wob
./deploy-partial-feature.sh wob true true
../remove-functions-kn/remove-partial-feature.sh wob
./deploy-partial-hotel.sh wob true true
../remove-functions-kn/remove-partial-hotel.sh wob
./deploy-partial-pipelined.sh wob true true
../remove-functions-kn/remove-partial-pipelined.sh wob
./deploy-partial-sentiment.sh wob true true
../remove-functions-kn/remove-partial-sentiment.sh wob
./deploy-partial-video.sh wob true true
../remove-functions-kn/remove-partial-video.sh wob
./deploy-partial-wage.sh wob true true
../remove-functions-kn/remove-partial-wage.sh wob
docker system prune -a -f
docker volume prune -f
