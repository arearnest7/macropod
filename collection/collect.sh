#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
benchmarks=("election" "hotel" "sentiment" "video")
unified=("election-unified" "hotel-unified" "sentiment-unified" "video-unified")
original=("election-gateway" "hotel-frontend" "sentiment-main" "video-streaming")
c=(8 16 20 40)
for i in {0..3}; do
	for concurrency in ${c[@]}; do
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		./collect-kn.sh $HOST ${benchmarks[$i]}-unified ${unified[$i]} $concurrency payloads/${benchmarks[$i]}.json 1000
		mv kn-${benchmarks[$i]}-unified.csv kn-unified-${benchmarks[$i]}-$concurrency.csv
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		./collect-kn.sh $HOST ${benchmarks[$i]}-original ${original[$i]} $concurrency payloads/${benchmarks[$i]}.json 1000
		mv kn-${benchmarks[$i]}-original.csv kn-original-${benchmarks[$i]}-$concurrency.csv
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		./collect-macropod.sh ${benchmarks[$i]} $concurrency payloads/${benchmarks[$i]}.json 1000
		mv macropod-${benchmarks[$i]}.csv macropod-${benchmarks[$i]}-$concurrency.csv
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
	done;
done;

benchmarks=("election" "hotel" "sentiment")
unified=("election-unified" "hotel-unified" "sentiment-unified")
original=("election-gateway" "hotel-frontend" "sentiment-main")
c=(80 100)
for i in {0..2}; do
        for concurrency in ${c[@]}; do
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                ./collect-kn.sh $HOST ${benchmarks[$i]}-unified ${unified[$i]} $concurrency payloads/${benchmarks[$i]}.json 1000
                mv kn-${benchmarks[$i]}-unified.csv kn-unified-${benchmarks[$i]}-$concurrency.csv
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                ./collect-kn.sh $HOST ${benchmarks[$i]}-original ${original[$i]} $concurrency payloads/${benchmarks[$i]}.json >
                mv kn-${benchmarks[$i]}-original.csv kn-original-${benchmarks[$i]}-$concurrency.csv
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                ./collect-macropod.sh ${benchmarks[$i]} $concurrency payloads/${benchmarks[$i]}.json 1000
                mv macropod-${benchmarks[$i]}.csv macropod-${benchmarks[$i]}-$concurrency.csv
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
        done;
done;

mkdir $DIR
chmod 777 $DIR
mv *.csv $DIR
mv *.out $DIR
mv cold-start $DIR
