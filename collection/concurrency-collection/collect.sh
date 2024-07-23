#!/bin/bash
HOST=${1:-127.0.0.1}
./collect-election.sh $HOST 1 500 >> election-1.out
./collect-election.sh $HOST 2 500 >> election-2.out
./collect-election.sh $HOST 4 500 >> election-4.out
./collect-election.sh $HOST 8 500 >> election-8.out
./collect-election.sh $HOST 16 500 >> election-16.out
./collect-election.sh $HOST 32 500 >> election-32.out
./collect-election.sh $HOST 64 500 >> election-64.out
./collect-election.sh $HOST 128 500 >> election-128.out
./collect-election.sh $HOST 256 500 >> election-256.out

#./collect-feature.sh $HOST 1 500 >> feature-1.out
#./collect-feature.sh $HOST 2 500 >> feature-2.out
#./collect-feature.sh $HOST 4 500 >> feature-4.out
#./collect-feature.sh $HOST 8 500 >> feature-8.out
#./collect-feature.sh $HOST 16 500 >> feature-16.out
#./collect-feature.sh $HOST 32 500 >> feature-32.out
#./collect-feature.sh $HOST 64 500 >> feature-64.out
#./collect-feature.sh $HOST 128 500 >> feature-128.out
#./collect-feature.sh $HOST 256 500 >> feature-256.out

./collect-hotel.sh $HOST 1 500 >> hotel-1.out
./collect-hotel.sh $HOST 2 500 >> hotel-2.out
./collect-hotel.sh $HOST 4 500 >> hotel-4.out
./collect-hotel.sh $HOST 8 500 >> hotel-8.out
./collect-hotel.sh $HOST 16 500 >> hotel-16.out
./collect-hotel.sh $HOST 32 500 >> hotel-32.out
./collect-hotel.sh $HOST 64 500 >> hotel-64.out
./collect-hotel.sh $HOST 128 500 >> hotel-128.out
./collect-hotel.sh $HOST 256 500 >> hotel-256.out

#./collect-pipelined.sh $HOST 1 500 >> pipelined-1.out
#./collect-pipelined.sh $HOST 2 500 >> pipelined-2.out
#./collect-pipelined.sh $HOST 4 500 >> pipelined-4.out
#./collect-pipelined.sh $HOST 8 500 >> pipelined-8.out
#./collect-pipelined.sh $HOST 16 500 >> pipelined-16.out
#./collect-pipelined.sh $HOST 32 500 >> pipelined-32.out
#./collect-pipelined.sh $HOST 64 500 >> pipelined-64.out
#./collect-pipelined.sh $HOST 128 500 >> pipelined-128.out
#./collect-pipelined.sh $HOST 256 500 >> pipelined-256.out

./collect-sentiment.sh $HOST 1 500 >> sentiment-1.out
./collect-sentiment.sh $HOST 2 500 >> sentiment-2.out
./collect-sentiment.sh $HOST 4 500 >> sentiment-4.out
./collect-sentiment.sh $HOST 8 500 >> sentiment-8.out
./collect-sentiment.sh $HOST 16 500 >> sentiment-16.out
./collect-sentiment.sh $HOST 32 500 >> sentiment-32.out
./collect-sentiment.sh $HOST 64 500 >> sentiment-64.out
./collect-sentiment.sh $HOST 128 500 >> sentiment-128.out
./collect-sentiment.sh $HOST 256 500 >> sentiment-256.out

./collect-video.sh $HOST 1 500 >> video-1.out
./collect-video.sh $HOST 2 500 >> video-2.out
./collect-video.sh $HOST 4 500 >> video-4.out
./collect-video.sh $HOST 8 500 >> video-8.out
./collect-video.sh $HOST 16 500 >> video-16.out
./collect-video.sh $HOST 32 500 >> video-32.out
./collect-video.sh $HOST 64 500 >> video-64.out
./collect-video.sh $HOST 128 500 >> video-128.out
./collect-video.sh $HOST 256 500 >> video-256.out

#./collect-wage.sh $HOST 1 500 >> wage-1.out
#./collect-wage.sh $HOST 2 500 >> wage-2.out
#./collect-wage.sh $HOST 4 500 >> wage-4.out
#./collect-wage.sh $HOST 8 500 >> wage-8.out
#./collect-wage.sh $HOST 16 500 >> wage-16.out
#./collect-wage.sh $HOST 32 500 >> wage-32.out
#./collect-wage.sh $HOST 64 500 >> wage-64.out
#./collect-wage.sh $HOST 128 500 >> wage-128.out
#./collect-wage.sh $HOST 256 500 >> wage-256.out
