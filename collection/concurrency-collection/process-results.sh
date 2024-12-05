#!/bin/bash
DIR=${1:-results}
METRIC_FILES=${2:-"$DIR/192.168.56.21-metrics.csv $DIR/192.168.56.22-metrics.csv $DIR/192.168.56.23-metrics.csv $DIR/192.168.56.24-metrics.csv"}
METRIC_RESULTS="$DIR/metrics-results.csv"
STATS_RESULTS="$DIR/stats-results.csv"

#TODO
LATENCY_RESULTS=()
LATENCY_OUT=()
benchmarks=("election" "hotel" "sentiment" "video")
c=(8 16 20 40)
for concurrency in ${c[@]}; do
	LATENCY_FILES=""
	LATENCY_OUT[$concurrency]=""
	for benchmark in ${benchmarks[@]}; do
		#LATENCY_FILES="${LATENCY_FILES} kn-full-${benchmark}-${concurrency}.csv"
		#LATENCY_FILES="${LATENCY_FILES} kn-original-${benchmark}-${concurrency}.csv"
		#LATENCY_FILES="${LATENCY_FILES} kn-partial-${benchmark}-${concurrency}.csv"
		LATENCY_FILES="${LATENCY_FILES} macropod-${benchmark}-${concurrency}.csv"
		if [[ "${LATENCY_OUT[$concurrency]}" == "" ]]; then
			LATENCY_OUT[$concurrency]="${benchmark}-${concurrency}.out"
		else
			LATENCY_OUT[$concurrency]="${LATENCY_OUT[$concurrency]};${benchmark}-${concurrency}.out"
		fi
	done;
	LATENCY_RESULTS[$concurrency]="${DIR}/latency-results-${concurrency}.csv"
	go run ../../tools/collection/processing-scripts/process-latency.go 1000 ${LATENCY_RESULTS[$concurrency]} $DIR $LATENCY_FILES
done;

benchmarks=("election" "hotel" "sentiment")
c=(80 100)
for concurrency in ${c[@]}; do
        LATENCY_FILES=""
	LATENCY_OUT[$concurrency]=""
        for benchmark in ${benchmarks[@]}; do
                #LATENCY_FILES="${LATENCY_FILES} kn-full-${benchmark}-${concurrency}.csv"
		#LATENCY_FILES="${LATENCY_FILES} kn-original-${benchmark}-${concurrency}.csv"
		#LATENCY_FILES="${LATENCY_FILES} kn-partial-${benchmark}-${concurrency}.csv"
		LATENCY_FILES="${LATENCY_FILES} macropod-${benchmark}-${concurrency}.csv"
		if [[ "${LATENCY_OUT[$concurrency]}" == "" ]]; then
                        LATENCY_OUT[$concurrency]="${benchmark}-${concurrency}.out"
                else
                        LATENCY_OUT[$concurrency]="${LATENCY_OUT[$concurrency]};${benchmark}-${concurrency}.out"
                fi
        done;
	LATENCY_RESULTS[$concurrency]="${DIR}/latency-results-${concurrency}.csv"
        go run ../../tools/collection/processing-scripts/process-latency.go 1000 ${LATENCY_RESULTS[$concurrency]} $DIR $LATENCY_FILES
done;

go run ../../tools/collection/processing-scripts/process-metrics-multi.go $METRIC_RESULTS $METRIC_FILES

c=(8 16 20 40 80 100)
LATENCY_STR=""
for concurrency in ${c[@]}; do
	LATENCY_STR="${LATENCY_STR} ${LATENCY_RESULTS[$concurrency]};${LATENCY_OUT[$concurrency]}"
done;
go run ../../tools/collection/processing-scripts/process-stats-multi.go $STATS_RESULTS $METRIC_RESULTS $DIR $LATENCY_STR
