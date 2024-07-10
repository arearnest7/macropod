#!/bin/bash
HOST=${1:-127.0.0.1}
./collect-election.sh $HOST 1 5 >> election-1.out
./collect-election.sh $HOST 2 10 >> election-2.out
./collect-election.sh $HOST 4 20 >> election-4.out
./collect-election.sh $HOST 8 40 >> election-8.out
./collect-election.sh $HOST 16 80 >> election-16.out
./collect-election.sh $HOST 32 160 >> election-32.out
./collect-election.sh $HOST 64 320 >> election-64.out
./collect-election.sh $HOST 128 640 >> election-128.out
#./collect-election.sh $HOST 256 1280 >> election-256.out

./collect-feature.sh $HOST 1 5 >> feature-1.out
./collect-feature.sh $HOST 2 10 >> feature-2.out
./collect-feature.sh $HOST 4 20 >> feature-4.out
#./collect-feature.sh $HOST 8 40 >> feature-8.out
#./collect-feature.sh $HOST 16 80 >> feature-16.out
#./collect-feature.sh $HOST 32 160 >> feature-32.out
#./collect-feature.sh $HOST 64 320 >> feature-64.out
#./collect-feature.sh $HOST 128 640 >> feature-128.out
#./collect-feature.sh $HOST 256 1280 >> feature-256.out

./collect-hotel.sh $HOST 1 5 >> hotel-1.out
./collect-hotel.sh $HOST 2 10 >> hotel-2.out
./collect-hotel.sh $HOST 4 20 >> hotel-4.out
./collect-hotel.sh $HOST 8 40 >> hotel-8.out
./collect-hotel.sh $HOST 16 80 >> hotel-16.out
./collect-hotel.sh $HOST 32 160 >> hotel-32.out
./collect-hotel.sh $HOST 64 320 >> hotel-64.out
./collect-hotel.sh $HOST 128 640 >> hotel-128.out
#./collect-hotel.sh $HOST 256 1280 >> hotel-256.out

./collect-pipelined.sh $HOST 1 5 >> pipelined-1.out
./collect-pipelined.sh $HOST 2 10 >> pipelined-2.out
./collect-pipelined.sh $HOST 4 20 >> pipelined-4.out
./collect-pipelined.sh $HOST 8 40 >> pipelined-8.out
./collect-pipelined.sh $HOST 16 80 >> pipelined-16.out
./collect-pipelined.sh $HOST 32 160 >> pipelined-32.out
./collect-pipelined.sh $HOST 64 320 >> pipelined-64.out
./collect-pipelined.sh $HOST 128 640 >> pipelined-128.out
#./collect-pipelined.sh $HOST 256 1280 >> pipelined-256.out

./collect-sentiment.sh $HOST 1 5 >> sentiment-1.out
./collect-sentiment.sh $HOST 2 10 >> sentiment-2.out
./collect-sentiment.sh $HOST 4 20 >> sentiment-4.out
./collect-sentiment.sh $HOST 8 40 >> sentiment-8.out
./collect-sentiment.sh $HOST 16 80 >> sentiment-16.out
./collect-sentiment.sh $HOST 32 160 >> sentiment-32.out
./collect-sentiment.sh $HOST 64 320 >> sentiment-64.out
./collect-sentiment.sh $HOST 128 640 >> sentiment-128.out
#./collect-sentiment.sh $HOST 256 1280 >> sentiment-256.out

./collect-video.sh $HOST 1 5 >> video-1.out
./collect-video.sh $HOST 2 10 >> video-2.out
./collect-video.sh $HOST 4 20 >> video-4.out
./collect-video.sh $HOST 8 40 >> video-8.out
./collect-video.sh $HOST 16 80 >> video-16.out
./collect-video.sh $HOST 32 160 >> video-32.out
./collect-video.sh $HOST 64 320 >> video-64.out
./collect-video.sh $HOST 128 640 >> video-128.out
#./collect-video.sh $HOST 256 1280 >> video-256.out

#./collect-wage.sh $HOST 1 5 >> wage-1.out
#./collect-wage.sh $HOST 2 10 >> wage-2.out
#./collect-wage.sh $HOST 4 20 >> wage-4.out
#./collect-wage.sh $HOST 8 40 >> wage-8.out
#./collect-wage.sh $HOST 16 80 >> wage-16.out
#./collect-wage.sh $HOST 32 160 >> wage-32.out
#./collect-wage.sh $HOST 64 320 >> wage-64.out
#./collect-wage.sh $HOST 128 640 >> wage-128.out
#./collect-wage.sh $HOST 256 1280 >> wage-256.out
