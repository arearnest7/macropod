#!/bin/bash
GATEWAY=${1:-10.125.189.107}
../../tools/deploy-scripts/deploy-original-election.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/election.json -o csv -t 180 http://election-gateway.default.$GATEWAY.sslip.io | tail -n +2 >> original-election-cold-start.csv

../../tools/remove-scripts/remove-original-election.sh
../../tools/deploy-scripts/deploy-full-election.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/election.json -o csv -t 180 http://election-full.default.$GATEWAY.sslip.io | tail -n +2 >> full-election-cold-start.csv

../../tools/remove-scripts/remove-full-election.sh
../../tools/deploy-scripts/deploy-partial-election.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/election.json -o csv -t 180 http://election-gateway-vevp.default.$GATEWAY.sslip.io | tail -n +2 >> partial-election-cold-start.csv

../../tools/remove-scripts/remove-partial-election.sh
../../tools/deploy-scripts/deploy-original-video.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/video.json -o csv -t 180 http://video-streaming.default.$GATEWAY.sslip.io | tail -n +2 >> original-video-cold-start.csv

../../tools/remove-scripts/remove-original-video.sh
../../tools/deploy-scripts/deploy-full-video.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/video.json -o csv -t 180 http://video-full.default.$GATEWAY.sslip.io | tail -n +2 >> full-video-cold-start.csv

../../tools/remove-scripts/remove-full-video.sh
../../tools/deploy-scripts/deploy-partial-video.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/video.json -o csv -t 180 http://video-streaming-d.default.$GATEWAY.sslip.io | tail -n +2 >> partial-video-cold-start.csv

../../tools/remove-scripts/remove-partial-video.sh
../../tools/deploy-scripts/deploy-original-hotel.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/hotel.json -o csv -t 180 http://hotel-frontend.default.$GATEWAY.sslip.io | tail -n +2 >> original-hotel-cold-start.csv

../../tools/remove-scripts/remove-original-hotel.sh
../../tools/deploy-scripts/deploy-full-hotel.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/hotel.json -o csv -t 180 http://hotel-full.default.$GATEWAY.sslip.io | tail -n +2 >> full-hotel-cold-start.csv

../../tools/remove-scripts/remove-full-hotel.sh
../../tools/deploy-scripts/deploy-partial-hotel.sh
sleep 30
hey -n 10000 -c 100 -D ../../tools/payloads/hotel.json -o csv -t 180 http://hotel-frontend-spgr.default.$GATEWAY.sslip.io | tail -n +2 >> partial-hotel-cold-start.csv

../../tools/remove-scripts/remove-partial-hotel.sh
