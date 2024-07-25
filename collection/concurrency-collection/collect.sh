#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
./collect-election.sh $HOST 1 500 >> election-1.out
./collect-election.sh $HOST 2 500 >> election-2.out
./collect-election.sh $HOST 4 500 >> election-4.out
./collect-election.sh $HOST 8 500 >> election-8.out
./collect-election.sh $HOST 16 500 >> election-16.out
./collect-election.sh $HOST 20 500 >> election-20.out
./collect-election.sh $HOST 30 500 >> election-30.out
./collect-election.sh $HOST 40 500 >> election-40.out
./collect-election.sh $HOST 50 500 >> election-50.out
./collect-election.sh $HOST 60 500 >> election-60.out
./collect-election.sh $HOST 70 500 >> election-70.out
./collect-election.sh $HOST 80 500 >> election-80.out
./collect-election.sh $HOST 90 500 >> election-90.out
./collect-election.sh $HOST 100 500 >> election-100.out
./collect-election.sh $HOST 110 500 >> election-110.out
./collect-election.sh $HOST 120 500 >> election-120.out
./collect-election.sh $HOST 130 500 >> election-130.out
./collect-election.sh $HOST 140 500 >> election-140.out
./collect-election.sh $HOST 150 500 >> election-150.out
./collect-election.sh $HOST 160 500 >> election-160.out
./collect-election.sh $HOST 170 500 >> election-170.out
./collect-election.sh $HOST 180 500 >> election-180.out
./collect-election.sh $HOST 190 500 >> election-190.out
./collect-election.sh $HOST 200 500 >> election-200.out
./collect-election.sh $HOST 210 500 >> election-210.out
./collect-election.sh $HOST 220 500 >> election-220.out
./collect-election.sh $HOST 230 500 >> election-230.out
./collect-election.sh $HOST 240 500 >> election-240.out
./collect-election.sh $HOST 250 500 >> election-250.out

#./collect-feature.sh $HOST 1 500 >> feature-1.out
#./collect-feature.sh $HOST 2 500 >> feature-2.out
#./collect-feature.sh $HOST 4 500 >> feature-4.out
#./collect-feature.sh $HOST 8 500 >> feature-8.out

./collect-hotel.sh $HOST 1 500 >> hotel-1.out
./collect-hotel.sh $HOST 2 500 >> hotel-2.out
./collect-hotel.sh $HOST 4 500 >> hotel-4.out
./collect-hotel.sh $HOST 8 500 >> hotel-8.out
./collect-hotel.sh $HOST 16 500 >> hotel-16.out
./collect-hotel.sh $HOST 20 500 >> hotel-20.out
./collect-hotel.sh $HOST 30 500 >> hotel-30.out
./collect-hotel.sh $HOST 40 500 >> hotel-40.out
./collect-hotel.sh $HOST 50 500 >> hotel-50.out
./collect-hotel.sh $HOST 60 500 >> hotel-60.out
./collect-hotel.sh $HOST 70 500 >> hotel-70.out
./collect-hotel.sh $HOST 80 500 >> hotel-80.out
./collect-hotel.sh $HOST 90 500 >> hotel-90.out
./collect-hotel.sh $HOST 100 500 >> hotel-100.out
./collect-hotel.sh $HOST 110 500 >> hotel-110.out
./collect-hotel.sh $HOST 120 500 >> hotel-120.out
./collect-hotel.sh $HOST 130 500 >> hotel-130.out
./collect-hotel.sh $HOST 140 500 >> hotel-140.out
./collect-hotel.sh $HOST 150 500 >> hotel-150.out
./collect-hotel.sh $HOST 160 500 >> hotel-160.out
./collect-hotel.sh $HOST 170 500 >> hotel-170.out
./collect-hotel.sh $HOST 180 500 >> hotel-180.out
./collect-hotel.sh $HOST 190 500 >> hotel-190.out
./collect-hotel.sh $HOST 200 500 >> hotel-200.out
./collect-hotel.sh $HOST 210 500 >> hotel-210.out
./collect-hotel.sh $HOST 220 500 >> hotel-220.out
./collect-hotel.sh $HOST 230 500 >> hotel-230.out
./collect-hotel.sh $HOST 240 500 >> hotel-240.out
./collect-hotel.sh $HOST 250 500 >> hotel-250.out

#./collect-pipelined.sh $HOST 1 500 >> pipelined-1.out
#./collect-pipelined.sh $HOST 2 500 >> pipelined-2.out
#./collect-pipelined.sh $HOST 4 500 >> pipelined-4.out
#./collect-pipelined.sh $HOST 8 500 >> pipelined-8.out
#./collect-pipelined.sh $HOST 16 500 >> pipelined-16.out

./collect-sentiment.sh $HOST 1 500 >> sentiment-1.out
./collect-sentiment.sh $HOST 2 500 >> sentiment-2.out
./collect-sentiment.sh $HOST 4 500 >> sentiment-4.out
./collect-sentiment.sh $HOST 8 500 >> sentiment-8.out
./collect-sentiment.sh $HOST 16 500 >> sentiment-16.out
./collect-sentiment.sh $HOST 20 500 >> sentiment-20.out
./collect-sentiment.sh $HOST 30 500 >> sentiment-30.out
./collect-sentiment.sh $HOST 40 500 >> sentiment-40.out
./collect-sentiment.sh $HOST 50 500 >> sentiment-50.out
./collect-sentiment.sh $HOST 60 500 >> sentiment-60.out
./collect-sentiment.sh $HOST 70 500 >> sentiment-70.out
./collect-sentiment.sh $HOST 80 500 >> sentiment-80.out
./collect-sentiment.sh $HOST 90 500 >> sentiment-90.out
./collect-sentiment.sh $HOST 100 500 >> sentiment-100.out
./collect-sentiment.sh $HOST 110 500 >> sentiment-110.out
./collect-sentiment.sh $HOST 120 500 >> sentiment-120.out
./collect-sentiment.sh $HOST 130 500 >> sentiment-130.out
./collect-sentiment.sh $HOST 140 500 >> sentiment-140.out
./collect-sentiment.sh $HOST 150 500 >> sentiment-150.out
./collect-sentiment.sh $HOST 160 500 >> sentiment-160.out
./collect-sentiment.sh $HOST 170 500 >> sentiment-170.out
./collect-sentiment.sh $HOST 180 500 >> sentiment-180.out
./collect-sentiment.sh $HOST 190 500 >> sentiment-190.out
./collect-sentiment.sh $HOST 200 500 >> sentiment-200.out
./collect-sentiment.sh $HOST 210 500 >> sentiment-210.out
./collect-sentiment.sh $HOST 220 500 >> sentiment-220.out
./collect-sentiment.sh $HOST 230 500 >> sentiment-230.out
./collect-sentiment.sh $HOST 240 500 >> sentiment-240.out
./collect-sentiment.sh $HOST 250 500 >> sentiment-250.out

./collect-video.sh $HOST 1 500 >> video-1.out
./collect-video.sh $HOST 2 500 >> video-2.out
./collect-video.sh $HOST 4 500 >> video-4.out
./collect-video.sh $HOST 8 500 >> video-8.out

#./collect-wage.sh $HOST 1 500 >> wage-1.out
#./collect-wage.sh $HOST 2 500 >> wage-2.out
#./collect-wage.sh $HOST 4 500 >> wage-4.out
#./collect-wage.sh $HOST 8 500 >> wage-8.out
#./collect-wage.sh $HOST 16 500 >> wage-16.out

mkdir $DIR
mv *.csv $DIR
mv *.out $DIR
mv cold-start $DIR
