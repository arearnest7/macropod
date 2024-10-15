#!/bin/bash
DIR=${1:-.}
METRIC_RESULTS=${2:-$DIR/metrics-results.csv}
STATS_RESULTS=${3:-$DIR/stats-results.csv}
METRIC_FILES=${4:-"$DIR/node-1-metrics.csv $DIR/node-2-metrics.csv $DIR/node-3-metrics.csv $DIR/node-4-metrics.csv $DIR/node-5-metrics.csv"}
LATENCY_FILES_1="kn-full-election-1.csv kn-original-election-1.csv kn-partial-election-1.csv kn-full-hotel-1.csv kn-original-hotel-1.csv kn-partial-hotel-1.csv kn-full-sentiment-1.csv kn-original-sentiment-1.csv kn-partial-sentiment-1.csv kn-full-video-1.csv kn-original-video-1.csv kn-partial-video-1.csv"
LATENCY_OUT_1="election-1.out;hotel-1.out;sentiment-1.out;video-1.out"
LATENCY_RESULTS_1="$DIR/latency-results-1.csv"
LATENCY_FILES_2="kn-full-election-2.csv kn-original-election-2.csv kn-partial-election-2.csv kn-full-hotel-2.csv kn-original-hotel-2.csv kn-partial-hotel-2.csv kn-full-sentiment-2.csv kn-original-sentiment-2.csv kn-partial-sentiment-2.csv kn-full-video-2.csv kn-original-video-2.csv kn-partial-video-2.csv"
LATENCY_OUT_2="election-2.out;hotel-2.out;sentiment-2.out;video-2.out"
LATENCY_RESULTS_2="$DIR/latency-results-2.csv"
LATENCY_FILES_4="kn-full-election-4.csv kn-original-election-4.csv kn-partial-election-4.csv kn-full-hotel-4.csv kn-original-hotel-4.csv kn-partial-hotel-4.csv kn-full-sentiment-4.csv kn-original-sentiment-4.csv kn-partial-sentiment-4.csv kn-full-video-4.csv kn-original-video-4.csv kn-partial-video-4.csv"
LATENCY_OUT_4="election-4.out;hotel-4.out;sentiment-4.out;video-4.out"
LATENCY_RESULTS_4="$DIR/latency-results-4.csv"
LATENCY_FILES_8="kn-full-election-8.csv kn-original-election-8.csv kn-partial-election-8.csv kn-full-hotel-8.csv kn-original-hotel-8.csv kn-partial-hotel-8.csv kn-full-sentiment-8.csv kn-original-sentiment-8.csv kn-partial-sentiment-8.csv kn-full-video-8.csv kn-original-video-8.csv kn-partial-video-8.csv"
LATENCY_OUT_8="election-8.out;hotel-8.out;sentiment-8.out;video-8.out"
LATENCY_RESULTS_8="$DIR/latency-results-8.csv"
LATENCY_FILES_16="kn-full-election-16.csv kn-original-election-16.csv kn-partial-election-16.csv kn-full-hotel-16.csv kn-original-hotel-16.csv kn-partial-hotel-16.csv kn-full-sentiment-16.csv kn-original-sentiment-16.csv kn-partial-sentiment-16.csv kn-full-video-16.csv kn-original-video-16.csv kn-partial-video-16.csv"
LATENCY_OUT_16="election-16.out;hotel-16.out;sentiment-16.out;video-16.out"
LATENCY_RESULTS_16="$DIR/latency-results-16.csv"
LATENCY_FILES_20="kn-full-election-20.csv kn-original-election-20.csv kn-partial-election-20.csv kn-full-hotel-20.csv kn-original-hotel-20.csv kn-partial-hotel-20.csv kn-full-sentiment-20.csv kn-original-sentiment-20.csv kn-partial-sentiment-20.csv kn-full-video-20.csv kn-original-video-20.csv kn-partial-video-20.csv"
LATENCY_OUT_20="election-20.out;hotel-20.out;sentiment-20.out;video-20.out"
LATENCY_RESULTS_20="$DIR/latency-results-20.csv"
LATENCY_FILES_40="kn-full-election-40.csv kn-original-election-40.csv kn-partial-election-40.csv kn-full-hotel-40.csv kn-original-hotel-40.csv kn-partial-hotel-40.csv kn-full-sentiment-40.csv kn-original-sentiment-40.csv kn-partial-sentiment-40.csv kn-full-video-40.csv kn-original-video-40.csv kn-partial-video-40.csv"
LATENCY_OUT_40="election-40.out;hotel-40.out;sentiment-40.out;video-40.out"
LATENCY_RESULTS_40="$DIR/latency-results-40.csv"
LATENCY_FILES_80="kn-full-election-80.csv kn-original-election-80.csv kn-partial-election-80.csv kn-full-hotel-80.csv kn-original-hotel-80.csv kn-partial-hotel-80.csv kn-full-sentiment-80.csv kn-original-sentiment-80.csv kn-partial-sentiment-80.csv"
LATENCY_OUT_80="election-80.out;hotel-80.out;sentiment-80.out"
LATENCY_RESULTS_80="$DIR/latency-results-80.csv"
LATENCY_FILES_100="kn-full-election-100.csv kn-original-election-100.csv kn-partial-election-100.csv kn-full-hotel-100.csv kn-original-hotel-100.csv kn-partial-hotel-100.csv kn-full-sentiment-100.csv kn-original-sentiment-100.csv kn-partial-sentiment-100.csv"
LATENCY_OUT_100="election-100.out;hotel-100.out;sentiment-100.out"
LATENCY_RESULTS_100="$DIR/latency-results-100.csv"
LATENCY_FILES_150="kn-full-election-150.csv kn-original-election-150.csv kn-partial-election-150.csv kn-full-hotel-150.csv kn-original-hotel-150.csv kn-partial-hotel-150.csv kn-full-sentiment-150.csv kn-original-sentiment-150.csv kn-partial-sentiment-150.csv"
LATENCY_OUT_150="election-150.out;hotel-150.out;sentiment-150.out"
LATENCY_RESULTS_150="$DIR/latency-results-150.csv"
LATENCY_FILES_200="kn-full-election-200.csv kn-original-election-200.csv kn-partial-election-200.csv kn-full-hotel-200.csv kn-original-hotel-200.csv kn-partial-hotel-200.csv kn-full-sentiment-200.csv kn-original-sentiment-200.csv kn-partial-sentiment-200.csv"
LATENCY_OUT_200="election-200.out;hotel-200.out;sentiment-200.out"
LATENCY_RESULTS_200="$DIR/latency-results-200.csv"
LATENCY_FILES_250="kn-full-election-250.csv kn-original-election-250.csv kn-partial-election-250.csv kn-full-hotel-250.csv kn-original-hotel-250.csv kn-partial-hotel-250.csv kn-full-sentiment-250.csv kn-original-sentiment-250.csv kn-partial-sentiment-250.csv"
LATENCY_OUT_250="election-250.out;hotel-250.out;sentiment-250.out"
LATENCY_RESULTS_250="$DIR/latency-results-250.csv"
LATENCY_FILES_300="kn-full-election-300.csv kn-original-election-300.csv kn-partial-election-300.csv kn-full-hotel-300.csv kn-original-hotel-300.csv kn-partial-hotel-300.csv kn-full-sentiment-300.csv kn-original-sentiment-300.csv kn-partial-sentiment-300.csv"
LATENCY_OUT_300="election-300.out;hotel-300.out;sentiment-300.out"
LATENCY_RESULTS_300="$DIR/latency-results-300.csv"
LATENCY_FILES_350="kn-full-election-350.csv kn-original-election-350.csv kn-partial-election-350.csv kn-full-hotel-350.csv kn-original-hotel-350.csv kn-partial-hotel-350.csv kn-full-sentiment-350.csv kn-original-sentiment-350.csv kn-partial-sentiment-350.csv"
LATENCY_OUT_350="election-350.out;hotel-350.out;sentiment-350.out"
LATENCY_RESULTS_350="$DIR/latency-results-350.csv"
LATENCY_FILES_400="kn-full-election-400.csv kn-original-election-400.csv kn-partial-election-400.csv kn-full-hotel-400.csv kn-original-hotel-400.csv kn-partial-hotel-400.csv kn-full-sentiment-400.csv kn-original-sentiment-400.csv kn-partial-sentiment-400.csv"
LATENCY_OUT_400="election-400.out;hotel-400.out;sentiment-400.out"
LATENCY_RESULTS_400="$DIR/latency-results-400.csv"
LATENCY_FILES_450="kn-full-election-450.csv kn-original-election-450.csv kn-partial-election-450.csv kn-full-hotel-450.csv kn-original-hotel-450.csv kn-partial-hotel-450.csv kn-full-sentiment-450.csv kn-original-sentiment-450.csv kn-partial-sentiment-450.csv"
LATENCY_OUT_450="election-450.out;hotel-450.out;sentiment-450.out"
LATENCY_RESULTS_450="$DIR/latency-results-450.csv"
LATENCY_FILES_500="kn-full-election-500.csv kn-original-election-500.csv kn-partial-election-500.csv kn-full-hotel-500.csv kn-original-hotel-500.csv kn-partial-hotel-500.csv kn-full-sentiment-500.csv kn-original-sentiment-500.csv kn-partial-sentiment-500.csv"
LATENCY_OUT_500="election-500.out;hotel-500.out;sentiment-500.out"
LATENCY_RESULTS_500="$DIR/latency-results-500.csv"
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_1 $DIR $LATENCY_FILES_1
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_2 $DIR $LATENCY_FILES_2
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_4 $DIR $LATENCY_FILES_4
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_8 $DIR $LATENCY_FILES_8
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_16 $DIR $LATENCY_FILES_16
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_20 $DIR $LATENCY_FILES_20
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_40 $DIR $LATENCY_FILES_40
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_80 $DIR $LATENCY_FILES_80
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_100 $DIR $LATENCY_FILES_100
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_150 $DIR $LATENCY_FILES_150
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_200 $DIR $LATENCY_FILES_200
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_250 $DIR $LATENCY_FILES_250
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_300 $DIR $LATENCY_FILES_300
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_350 $DIR $LATENCY_FILES_350
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_400 $DIR $LATENCY_FILES_400
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_450 $DIR $LATENCY_FILES_450
go run ../../tools/collection/processing-scripts/process-latency.go 1000 $LATENCY_RESULTS_500 $DIR $LATENCY_FILES_500

go run ../../tools/collection/processing-scripts/process-metrics-multi.go $METRIC_RESULTS $METRIC_FILES

go run ../../tools/collection/processing-scripts/process-stats-multi.go $STATS_RESULTS $METRIC_RESULTS $DIR "$LATENCY_RESULTS_1;$LATENCY_OUT_1" "$LATENCY_RESULTS_2;$LATENCY_OUT_2" "$LATENCY_RESULTS_4;$LATENCY_OUT_4" "$LATENCY_RESULTS_8;$LATENCY_OUT_8" "$LATENCY_RESULTS_16;$LATENCY_OUT_16" "$LATENCY_RESULTS_20;$LATENCY_OUT_20" "$LATENCY_RESULTS_40;$LATENCY_OUT_40" "$LATENCY_RESULTS_80;$LATENCY_OUT_80" "$LATENCY_RESULTS_100;$LATENCY_OUT_100" "$LATENCY_RESULTS_150;$LATENCY_OUT_150" "$LATENCY_RESULTS_200;$LATENCY_OUT_200" "$LATENCY_RESULTS_250;$LATENCY_OUT_250" "$LATENCY_RESULTS_300;$LATENCY_OUT_300" "$LATENCY_RESULTS_350;$LATENCY_OUT_350" "$LATENCY_RESULTS_400;$LATENCY_OUT_400" "$LATENCY_RESULTS_450;$LATENCY_OUT_450" "$LATENCY_RESULTS_500;$LATENCY_OUT_500"
