#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
./collect-election.sh $HOST 1 1000 >> election-1.out
./collect-election.sh $HOST 2 1000 >> election-2.out
./collect-election.sh $HOST 4 1000 >> election-4.out
./collect-election.sh $HOST 8 1000 >> election-8.out
./collect-election.sh $HOST 16 1000 >> election-16.out
./collect-election.sh $HOST 20 1000 >> election-20.out
./collect-election.sh $HOST 40 1000 >> election-40.out
./collect-election.sh $HOST 80 1000 >> election-80.out
./collect-election.sh $HOST 100 1000 >> election-100.out
./collect-election.sh $HOST 150 1000 >> election-150.out
./collect-election.sh $HOST 200 1000 >> election-200.out
./collect-election.sh $HOST 250 1000 >> election-250.out
./collect-election.sh $HOST 300 1000 >> election-300.out
./collect-election.sh $HOST 350 1000 >> election-350.out
./collect-election.sh $HOST 400 1000 >> election-400.out
./collect-election.sh $HOST 450 1000 >> election-450.out
./collect-election.sh $HOST 500 1000 >> election-500.out

./collect-hotel.sh $HOST 1 1000 >> hotel-1.out
./collect-hotel.sh $HOST 2 1000 >> hotel-2.out
./collect-hotel.sh $HOST 4 1000 >> hotel-4.out
./collect-hotel.sh $HOST 8 1000 >> hotel-8.out
./collect-hotel.sh $HOST 16 1000 >> hotel-16.out
./collect-hotel.sh $HOST 20 1000 >> hotel-20.out
./collect-hotel.sh $HOST 40 1000 >> hotel-40.out
./collect-hotel.sh $HOST 80 1000 >> hotel-80.out
./collect-hotel.sh $HOST 100 1000 >> hotel-100.out
./collect-hotel.sh $HOST 150 1000 >> hotel-150.out
./collect-hotel.sh $HOST 200 1000 >> hotel-200.out
./collect-hotel.sh $HOST 250 1000 >> hotel-250.out
./collect-hotel.sh $HOST 300 1000 >> hotel-300.out
./collect-hotel.sh $HOST 350 1000 >> hotel-350.out
./collect-hotel.sh $HOST 400 1000 >> hotel-400.out
./collect-hotel.sh $HOST 450 1000 >> hotel-450.out
./collect-hotel.sh $HOST 500 1000 >> hotel-500.out

./collect-sentiment.sh $HOST 1 1000 >> sentiment-1.out
./collect-sentiment.sh $HOST 2 1000 >> sentiment-2.out
./collect-sentiment.sh $HOST 4 1000 >> sentiment-4.out
./collect-sentiment.sh $HOST 8 1000 >> sentiment-8.out
./collect-sentiment.sh $HOST 16 1000 >> sentiment-16.out
./collect-sentiment.sh $HOST 20 1000 >> sentiment-20.out
./collect-sentiment.sh $HOST 40 1000 >> sentiment-40.out
./collect-sentiment.sh $HOST 80 1000 >> sentiment-80.out
./collect-sentiment.sh $HOST 100 1000 >> sentiment-100.out
./collect-sentiment.sh $HOST 150 1000 >> sentiment-150.out
./collect-sentiment.sh $HOST 200 1000 >> sentiment-200.out
./collect-sentiment.sh $HOST 250 1000 >> sentiment-250.out
./collect-sentiment.sh $HOST 300 1000 >> sentiment-300.out
./collect-sentiment.sh $HOST 350 1000 >> sentiment-350.out
./collect-sentiment.sh $HOST 400 1000 >> sentiment-400.out
./collect-sentiment.sh $HOST 450 1000 >> sentiment-450.out
./collect-sentiment.sh $HOST 500 1000 >> sentiment-500.out

./collect-video.sh $HOST 1 1000 >> video-1.out
./collect-video.sh $HOST 2 1000 >> video-2.out
./collect-video.sh $HOST 4 1000 >> video-4.out
./collect-video.sh $HOST 8 1000 >> video-8.out
./collect-video.sh $HOST 16 1000 >> video-16.out
./collect-video.sh $HOST 20 1000 >> video-20.out
./collect-video.sh $HOST 40 1000 >> video-40.out

mkdir $DIR
chmod 777 $DIR
mv *.csv $DIR
mv *.out $DIR
mv cold-start $DIR
