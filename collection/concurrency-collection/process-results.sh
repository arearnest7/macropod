#!/bin/bash
DIR=${1:-.}
METRIC_RESULTS=${2:-$DIR/metrics-results.csv}
STATS_RESULTS=${3:-$DIR/stats-results.csv}
METRIC_FILES=${4:-"$DIR/node-1-metrics.csv $DIR/node-2-metrics.csv $DIR/node-3-metrics.csv $DIR/node-4-metrics.csv $DIR/node-5-metrics.csv"}
LATENCY_FILES_1=${5:-"kn-full-election-1.csv kn-original-election-1.csv kn-partial-election-1.csv kn-full-hotel-1.csv kn-original-hotel-1.csv kn-partial-hotel-1.csv kn-full-sentiment-1.csv kn-original-sentiment-1.csv kn-partial-sentiment-1.csv kn-full-video-1.csv kn-original-video-1.csv kn-partial-video-1.csv"}
LATENCY_OUT_1=${6:-"election-1.out;feature-1.out;hotel-1.out;pipelined-1.out;sentiment-1.out;video-1.out"}
LATENCY_RESULTS_1=${7:-$DIR/latency-results-1.csv}
LATENCY_FILES_2=${8:-"kn-full-election-2.csv kn-original-election-2.csv kn-partial-election-2.csv kn-full-hotel-2.csv kn-original-hotel-2.csv kn-partial-hotel-2.csv kn-full-sentiment-2.csv kn-original-sentiment-2.csv kn-partial-sentiment-2.csv kn-full-video-2.csv kn-original-video-2.csv kn-partial-video-2.csv"}
LATENCY_OUT_2=${9:-"election-2.out;feature-2.out;hotel-2.out;pipelined-2.out;sentiment-2.out;video-2.out"}
LATENCY_RESULTS_2=${10:-$DIR/latency-results-2.csv}
LATENCY_FILES_4=${11:-"kn-full-election-4.csv kn-original-election-4.csv kn-partial-election-4.csv kn-full-hotel-4.csv kn-original-hotel-4.csv kn-partial-hotel-4.csv kn-full-sentiment-4.csv kn-original-sentiment-4.csv kn-partial-sentiment-4.csv kn-full-video-4.csv kn-original-video-4.csv kn-partial-video-4.csv"}
LATENCY_OUT_4=${12:-"election-4.out;feature-4.out;hotel-4.out;pipelined-4.out;sentiment-4.out;video-4.out"}
LATENCY_RESULTS_4=${13:-$DIR/latency-results-4.csv}
LATENCY_FILES_8=${14:-"kn-full-election-8.csv kn-original-election-8.csv kn-partial-election-8.csv kn-full-hotel-8.csv kn-original-hotel-8.csv kn-partial-hotel-8.csv kn-full-sentiment-8.csv kn-original-sentiment-8.csv kn-partial-sentiment-8.csv kn-full-video-8.csv kn-original-video-8.csv kn-partial-video-8.csv"}
LATENCY_OUT_8=${15:-"election-8.out;hotel-8.out;pipelined-8.out;sentiment-8.out;video-8.out"}
LATENCY_RESULTS_8=${16:-$DIR/latency-results-8.csv}
LATENCY_FILES_16=${17:-"kn-full-election-16.csv kn-original-election-16.csv kn-partial-election-16.csv kn-full-hotel-16.csv kn-original-hotel-16.csv kn-partial-hotel-16.csv kn-full-sentiment-16.csv kn-original-sentiment-16.csv kn-partial-sentiment-16.csv kn-full-video-16.csv kn-original-video-16.csv kn-partial-video-16.csv"}
LATENCY_OUT_16=${18:-"election-16.out;hotel-16.out;pipelined-16.out;sentiment-16.out;video-16.out"}
LATENCY_RESULTS_16=${19:-$DIR/latency-results-16.csv}
LATENCY_FILES_32=${20:-"kn-full-election-32.csv kn-original-election-32.csv kn-partial-election-32.csv kn-full-hotel-32.csv kn-original-hotel-32.csv kn-partial-hotel-32.csv kn-full-sentiment-32.csv kn-original-sentiment-32.csv kn-partial-sentiment-32.csv kn-full-video-32.csv kn-original-video-32.csv kn-partial-video-32.csv"}
LATENCY_OUT_32=${21:-"election-32.out;hotel-32.out;pipelined-32.out;sentiment-32.out;video-32.out"}
LATENCY_RESULTS_32=${22:-$DIR/latency-results-32.csv}
LATENCY_FILES_64=${23:-"kn-full-election-64.csv kn-original-election-64.csv kn-partial-election-64.csv kn-full-hotel-64.csv kn-original-hotel-64.csv kn-partial-hotel-64.csv kn-full-sentiment-64.csv kn-original-sentiment-64.csv kn-partial-sentiment-64.csv kn-full-video-64.csv kn-original-video-64.csv kn-partial-video-64.csv"}
LATENCY_OUT_64=${24:-"election-64.out;hotel-64.out;pipelined-64.out;sentiment-64.out;video-64.out"}
LATENCY_RESULTS_64=${25:-$DIR/latency-results-64.csv}
LATENCY_FILES_128=${26:-"kn-full-election-128.csv kn-original-election-128.csv kn-partial-election-128.csv kn-full-hotel-128.csv kn-original-hotel-128.csv kn-partial-hotel-128.csv kn-full-sentiment-128.csv kn-original-sentiment-128.csv kn-partial-sentiment-128.csv kn-full-video-128.csv kn-original-video-128.csv kn-partial-video-128.csv"}
LATENCY_OUT_128=${27:-"election-128.out;hotel-128.out;pipelined-128.out;sentiment-128.out;video-128.out"}
LATENCY_RESULTS_128=${28:-$DIR/latency-results-128.csv}
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_1 $DIR $LATENCY_FILES_1
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_2 $DIR $LATENCY_FILES_2
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_4 $DIR $LATENCY_FILES_4
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_8 $DIR $LATENCY_FILES_8
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_16 $DIR $LATENCY_FILES_16
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_32 $DIR $LATENCY_FILES_32
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_64 $DIR $LATENCY_FILES_64
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_128 $DIR $LATENCY_FILES_128

go run ../../tools/collection/processing-scripts/process-metrics-all.go $METRIC_RESULTS $METRIC_FILES

go run ../../tools/collection/processing-scripts/process-multi-stats.go $STATS_RESULTS $METRIC_RESULTS $DIR "$LATENCY_RESULTS_1;$LATENCY_OUT_1" "$LATENCY_RESULTS_2;$LATENCY_OUT_2" "$LATENCY_RESULTS_4;$LATENCY_OUT_4" "$LATENCY_RESULTS_8;$LATENCY_OUT_8" "$LATENCY_RESULTS_16;$LATENCY_OUT_16" "$LATENCY_RESULTS_32;$LATENCY_OUT_32" "$LATENCY_RESULTS_64;$LATENCY_OUT_64" "$LATENCY_RESULTS_128;$LATENCY_OUT_128"
