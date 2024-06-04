#!/bin/bash
kn func deploy --build=true --push=true --path ../../../benchmarks/micro/micro-rpc-a
kn func deploy --build=true --push=true --path ../../../benchmarks/micro/micro-rpc-b
kn func deploy --build=true --push=true --path ../../../benchmarks/micro/micro-rpc-a-b
kn func delete micro-rpc-a
kn func delete micro-rpc-b
kn func delete micro-rpc-a-b
../kn/deploy-full-election.sh kn true true
../remove-kn/remove-full-election.sh kn
../kn/deploy-full-feature.sh kn true true
../remove-kn/remove-full-feature.sh kn
../kn/deploy-full-hotel.sh kn true true
../remove-kn/remove-full-hotel.sh kn
../kn/deploy-full-pipelined.sh kn true true
../remove-kn/remove-full-pipelined.sh kn
../kn/deploy-full-sentiment.sh kn true true
../remove-kn/remove-full-sentiment.sh kn
../kn/deploy-full-video.sh kn true true
../remove-kn/remove-full-video.sh kn
../kn/deploy-full-wage.sh kn true true
../remove-kn/remove-full-wage.sh kn
../kn/deploy-original-election.sh kn true true
../remove-kn/remove-original-election.sh kn
../kn/deploy-original-feature.sh kn true true
../remove-kn/remove-original-feature.sh kn
../kn/deploy-original-hotel.sh kn true true
../remove-kn/remove-original-hotel.sh kn
../kn/deploy-original-pipelined.sh kn true true
../remove-kn/remove-original-pipelined.sh kn
../kn/deploy-original-sentiment.sh kn true true
../remove-kn/remove-original-sentiment.sh kn
../kn/deploy-original-video.sh kn true true
../remove-kn/remove-original-video.sh kn
../kn/deploy-original-wage.sh kn true true
../remove-kn/remove-original-wage.sh kn
../kn/deploy-partial-election.sh kn true true
../remove-kn/remove-partial-election.sh kn
../kn/deploy-partial-feature.sh kn true true
../remove-kn/remove-partial-feature.sh kn
../kn/deploy-partial-hotel.sh kn true true
../remove-kn/remove-partial-hotel.sh kn
../kn/deploy-partial-pipelined.sh kn true true
../remove-kn/remove-partial-pipelined.sh kn
../kn/deploy-partial-sentiment.sh kn true true
../remove-kn/remove-partial-sentiment.sh kn
../kn/deploy-partial-video.sh kn true true
../remove-kn/remove-partial-video.sh kn
../kn/deploy-partial-wage.sh kn true true
../remove-kn/remove-partial-wage.sh kn
../kn/deploy-full-election.sh wob true true
../remove-kn/remove-full-election.sh wob
../kn/deploy-full-feature.sh wob true true
../remove-kn/remove-full-feature.sh wob
../kn/deploy-full-hotel.sh wob true true
../remove-kn/remove-full-hotel.sh wob
../kn/deploy-full-pipelined.sh wob true true
../remove-kn/remove-full-pipelined.sh wob
../kn/deploy-full-sentiment.sh wob true true
../remove-kn/remove-full-sentiment.sh wob
../kn/deploy-full-video.sh wob true true
../remove-kn/remove-full-video.sh wob
../kn/deploy-full-wage.sh wob true true
../remove-kn/remove-full-wage.sh wob
../kn/deploy-original-election.sh wob true true
../remove-kn/remove-original-election.sh wob
../kn/deploy-original-feature.sh wob true true
../remove-kn/remove-original-feature.sh wob
../kn/deploy-original-hotel.sh wob true true
../remove-kn/remove-original-hotel.sh wob
../kn/deploy-original-pipelined.sh wob true true
../remove-kn/remove-original-pipelined.sh wob
../kn/deploy-original-sentiment.sh wob true true
../remove-kn/remove-original-sentiment.sh wob
../kn/deploy-original-video.sh wob true true
../remove-kn/remove-original-video.sh wob
../kn/deploy-original-wage.sh wob true true
../remove-kn/remove-original-wage.sh wob
../kn/deploy-partial-election.sh wob true true
../remove-kn/remove-partial-election.sh wob
../kn/deploy-partial-feature.sh wob true true
../remove-kn/remove-partial-feature.sh wob
../kn/deploy-partial-hotel.sh wob true true
../remove-kn/remove-partial-hotel.sh wob
../kn/deploy-partial-pipelined.sh wob true true
../remove-kn/remove-partial-pipelined.sh wob
../kn/deploy-partial-sentiment.sh wob true true
../remove-kn/remove-partial-sentiment.sh wob
../kn/deploy-partial-video.sh wob true true
../remove-kn/remove-partial-video.sh wob
../kn/deploy-partial-wage.sh wob true true
../remove-kn/remove-partial-wage.sh wob
