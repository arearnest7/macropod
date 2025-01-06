#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
benchmarks=("election" "hotel" "sentiment" "video")
full=("election-full" "hotel-full" "sentiment-full" "video-full")
original=("election-gateway" "hotel-frontend" "sentiment-main" "video-streaming")
partial=("election-gateway-vevp" "hotel-frontend-spgr" "sentiment-main-rcposc" "video-streaming-d")
c=(8 16 20 40)
for i in {0..1}; do
	for concurrency in ${c[@]}; do
		#date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		#./collect-kn.sh $HOST full-${benchmarks[$i]} ${full[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
		#mv kn-full-$benchmarks[$i].csv kn-full-$benchmarks[$i]-$concurrency.csv
		#date -u '+%F %H:%M:%S.%6N %Z' >> $benchmarks[$i]-$concurrency.out
		#./collect-kn.sh $HOST original-${benchmarks[$i]} ${original[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
		#mv kn-original-${benchmarks[$i]}.csv kn-original-${benchmarks[$i]}-$concurrency.csv
		#date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		#./collect-kn.sh $HOST partial-${benchmarks[$i]} ${partial[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
		#mv kn-partial-$benchmarks[$i].csv kn-partial-$benchmarks[$i]-$concurrency.csv
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
		./collect-macropod.sh ${benchmarks[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
		mv macropod-${benchmarks[$i]}.csv macropod-${benchmarks[$i]}-$concurrency.csv
		date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
	done;
done;

benchmarks=("election" "hotel" "sentiment")
full=("election-full" "hotel-full" "sentiment-full")
original=("election-gateway" "hotel-frontend" "sentiment-main")
partial=("election-gateway-vevp" "hotel-frontend-spgr" "sentiment-main-rcposc")
c=(80 100)
for i in {0..1}; do
        for concurrency in ${c[@]}; do
                #date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                #./collect-kn.sh $HOST full-${benchmarks[$i]} ${full[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
                #mv kn-full-$benchmarks[$i].csv kn-full-$benchmarks[$i]-$concurrency.csv
                #date -u '+%F %H:%M:%S.%6N %Z' >> $benchmarks[$i]-$concurrency.out
                #./collect-kn.sh $HOST original-${benchmarks[$i]} ${original[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json >
                #mv kn-original-${benchmarks[$i]}.csv kn-original-${benchmarks[$i]}-$concurrency.csv
                #date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                #./collect-kn.sh $HOST partial-${benchmarks[$i]} ${partial[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 10>
                #mv kn-partial-$benchmarks[$i].csv kn-partial-$benchmarks[$i]-$concurrency.csv
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
                ./collect-macropod.sh ${benchmarks[$i]} $concurrency ../../tools/payloads/${benchmarks[$i]}.json 1000
                mv macropod-${benchmarks[$i]}.csv macropod-${benchmarks[$i]}-$concurrency.csv
                date -u '+%F %H:%M:%S.%6N %Z' >> ${benchmarks[$i]}-$concurrency.out
        done;
done;

mkdir $DIR
chmod 777 $DIR
mv *.csv $DIR
mv *.out $DIR
mv cold-start $DIR
