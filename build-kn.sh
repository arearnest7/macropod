#!/bin/bash
./tools/deploy-functions-kn/deploy-full-election.sh kn true true
./tools/remove-functions-kn/remove-full-election.sh kn
./tools/deploy-functions-kn/deploy-full-feature.sh kn true true
./tools/remove-functions-kn/remove-full-feature.sh kn
./tools/deploy-functions-kn/deploy-full-hotel.sh kn true true
./tools/remove-functions-kn/remove-full-hotel.sh kn
./tools/deploy-functions-kn/deploy-full-pipelined.sh kn true true
./tools/remove-functions-kn/remove-full-pipelined.sh kn
./tools/deploy-functions-kn/deploy-full-sentiment.sh kn true true
./tools/remove-functions-kn/remove-full-sentiment.sh kn
./tools/deploy-functions-kn/deploy-full-video.sh kn true true
./tools/remove-functions-kn/remove-full-video.sh kn
./tools/deploy-functions-kn/deploy-full-wage.sh kn true true
./tools/remove-functions-kn/remove-full-wage.sh kn
sudo docker system prune -a -f
sudo docker volume prune -f
./tools/deploy-functions-kn/deploy-original-election.sh kn true true
./tools/remove-functions-kn/remove-original-election.sh kn
./tools/deploy-functions-kn/deploy-original-feature.sh kn true true
./tools/remove-functions-kn/remove-original-feature.sh kn
./tools/deploy-functions-kn/deploy-original-hotel.sh kn true true
./tools/remove-functions-kn/remove-original-hotel.sh kn
./tools/deploy-functions-kn/deploy-original-pipelined.sh kn true true
./tools/remove-functions-kn/remove-original-pipelined.sh kn
./tools/deploy-functions-kn/deploy-original-sentiment.sh kn true true
./tools/remove-functions-kn/remove-original-sentiment.sh kn
./tools/deploy-functions-kn/deploy-original-video.sh kn true true
./tools/remove-functions-kn/remove-original-video.sh kn
./tools/deploy-functions-kn/deploy-original-wage.sh kn true true
./tools/remove-functions-kn/remove-original-wage.sh kn
sudo docker system prune -a -f
sudo docker volume prune -f
./tools/deploy-functions-kn/deploy-partial-election.sh kn true true
./tools/remove-functions-kn/remove-partial-election.sh kn
./tools/deploy-functions-kn/deploy-partial-feature.sh kn true true
./tools/remove-functions-kn/remove-partial-feature.sh kn
./tools/deploy-functions-kn/deploy-partial-hotel.sh kn true true
./tools/remove-functions-kn/remove-partial-hotel.sh kn
./tools/deploy-functions-kn/deploy-partial-pipelined.sh kn true true
./tools/remove-functions-kn/remove-partial-pipelined.sh kn
./tools/deploy-functions-kn/deploy-partial-sentiment.sh kn true true
./tools/remove-functions-kn/remove-partial-sentiment.sh kn
./tools/deploy-functions-kn/deploy-partial-video.sh kn true true
./tools/remove-functions-kn/remove-partial-video.sh kn
./tools/deploy-functions-kn/deploy-partial-wage.sh kn true true
./tools/remove-functions-kn/remove-partial-wage.sh kn
sudo docker system prune -a -f
sudo docker volume prune -f
./tools/deploy-functions-kn/deploy-full-election.sh wob true true
./tools/remove-functions-kn/remove-full-election.sh wob
./tools/deploy-functions-kn/deploy-full-feature.sh wob true true
./tools/remove-functions-kn/remove-full-feature.sh wob
./tools/deploy-functions-kn/deploy-full-hotel.sh wob true true
./tools/remove-functions-kn/remove-full-hotel.sh wob
./tools/deploy-functions-kn/deploy-full-pipelined.sh wob true true
./tools/remove-functions-kn/remove-full-pipelined.sh wob
./tools/deploy-functions-kn/deploy-full-sentiment.sh wob true true
./tools/remove-functions-kn/remove-full-sentiment.sh wob
./tools/deploy-functions-kn/deploy-full-video.sh wob true true
./tools/remove-functions-kn/remove-full-video.sh wob
./tools/deploy-functions-kn/deploy-full-wage.sh wob true true
./tools/remove-functions-kn/remove-full-wage.sh wob
sudo docker system prune -a -f
sudo docker volume prune -f
./tools/deploy-functions-kn/deploy-original-election.sh wob true true
./tools/remove-functions-kn/remove-original-election.sh wob
./tools/deploy-functions-kn/deploy-original-feature.sh wob true true
./tools/remove-functions-kn/remove-original-feature.sh wob
./tools/deploy-functions-kn/deploy-original-hotel.sh wob true true
./tools/remove-functions-kn/remove-original-hotel.sh wob
./tools/deploy-functions-kn/deploy-original-pipelined.sh wob true true
./tools/remove-functions-kn/remove-original-pipelined.sh wob
./tools/deploy-functions-kn/deploy-original-sentiment.sh wob true true
./tools/remove-functions-kn/remove-original-sentiment.sh wob
./tools/deploy-functions-kn/deploy-original-video.sh wob true true
./tools/remove-functions-kn/remove-original-video.sh wob
./tools/deploy-functions-kn/deploy-original-wage.sh wob true true
./tools/remove-functions-kn/remove-original-wage.sh wob
sudo docker system prune -a -f
sudo docker volume prune -f
./tools/deploy-functions-kn/deploy-partial-election.sh wob true true
./tools/remove-functions-kn/remove-partial-election.sh wob
./tools/deploy-functions-kn/deploy-partial-feature.sh wob true true
./tools/remove-functions-kn/remove-partial-feature.sh wob
./tools/deploy-functions-kn/deploy-partial-hotel.sh wob true true
./tools/remove-functions-kn/remove-partial-hotel.sh wob
./tools/deploy-functions-kn/deploy-partial-pipelined.sh wob true true
./tools/remove-functions-kn/remove-partial-pipelined.sh wob
./tools/deploy-functions-kn/deploy-partial-sentiment.sh wob true true
./tools/remove-functions-kn/remove-partial-sentiment.sh wob
./tools/deploy-functions-kn/deploy-partial-video.sh wob true true
./tools/remove-functions-kn/remove-partial-video.sh wob
./tools/deploy-functions-kn/deploy-partial-wage.sh wbo true true
./tools/remove-functions-kn/remove-partial-wage.sh wob
sudo docker system prune -a -f
sudo docker volume prune -f
