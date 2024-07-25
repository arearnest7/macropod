#!/bin/bash
DIR=${1:-.}
METRIC_RESULTS=${2:-$DIR/metrics-results.csv}
STATS_RESULTS=${3:-$DIR/stats-results.csv}
METRIC_FILES=${4:-"$DIR/node-1-metrics.csv $DIR/node-2-metrics.csv $DIR/node-3-metrics.csv $DIR/node-4-metrics.csv $DIR/node-5-metrics.csv"}
LATENCY_FILES_1="kn-full-election-1.csv kn-original-election-1.csv kn-partial-election-1.csv kn-full-hotel-1.csv kn-original-hotel-1.csv kn-partial-hotel-1.csv kn-full-sentiment-1.csv kn-original-sentiment-1.csv kn-partial-sentiment-1.csv kn-full-video-1.csv kn-original-video-1.csv kn-partial-video-1.csv"
LATENCY_OUT_1="election-1.out;hotel-1.out;sentiment-1.out;video-1.out"
LATENCY_RESULTS_1="$DIR/latency-results-1.csv"
LATENCY_FILES_2="kn-full-election-2.csv kn-original-election-2.csv kn-partial-election-2.csv kn-full-hotel-2.csv kn-original-hotel-2.csv kn-partial-hotel-2.csv kn-full-sentiment-2.csv kn-original-sentiment-2.csv kn-partial-sentiment-2.csv kn-full-video-2.csv kn-original-video-2.csv kn-partial-video-2.csv"
LATENCY_OUT_2=-"election-2.out;hotel-2.out;sentiment-2.out;video-2.out"
LATENCY_RESULTS_2="$DIR/latency-results-2.csv"
LATENCY_FILES_4="kn-full-election-4.csv kn-original-election-4.csv kn-partial-election-4.csv kn-full-hotel-4.csv kn-original-hotel-4.csv kn-partial-hotel-4.csv kn-full-sentiment-4.csv kn-original-sentiment-4.csv kn-partial-sentiment-4.csv kn-full-video-4.csv kn-original-video-4.csv kn-partial-video-4.csv"
LATENCY_OUT_4="election-4.out;hotel-4.out;sentiment-4.out;video-4.out"
LATENCY_RESULTS_4="$DIR/latency-results-4.csv"
LATENCY_FILES_8="kn-full-election-8.csv kn-original-election-8.csv kn-partial-election-8.csv kn-full-hotel-8.csv kn-original-hotel-8.csv kn-partial-hotel-8.csv kn-full-sentiment-8.csv kn-original-sentiment-8.csv kn-partial-sentiment-8.csv kn-full-video-8.csv kn-original-video-8.csv kn-partial-video-8.csv"
LATENCY_OUT_8="election-8.out;hotel-8.out;sentiment-8.out;video-8.out"
LATENCY_RESULTS_8="$DIR/latency-results-8.csv"
LATENCY_FILES_16="kn-full-election-16.csv kn-original-election-16.csv kn-partial-election-16.csv kn-full-hotel-16.csv kn-original-hotel-16.csv kn-partial-hotel-16.csv kn-full-sentiment-16.csv kn-original-sentiment-16.csv kn-partial-sentiment-16.csv"
LATENCY_OUT_16="election-16.out;hotel-16.out;sentiment-16.out"
LATENCY_RESULTS_16="$DIR/latency-results-16.csv"
LATENCY_FILES_20="kn-full-election-20.csv kn-original-election-20.csv kn-partial-election-20.csv kn-full-hotel-20.csv kn-original-hotel-20.csv kn-partial-hotel-20.csv kn-full-sentiment-20.csv kn-original-sentiment-20.csv kn-partial-sentiment-20.csv"
LATENCY_OUT_20="election-20.out;hotel-20.out;sentiment-20.out"
LATENCY_RESULTS_20="$DIR/latency-results-20.csv"
LATENCY_FILES_30="kn-full-election-30.csv kn-original-election-30.csv kn-partial-election-30.csv kn-full-hotel-30.csv kn-original-hotel-30.csv kn-partial-hotel-30.csv kn-full-sentiment-30.csv kn-original-sentiment-30.csv kn-partial-sentiment-30.csv"
LATENCY_OUT_30="election-30.out;hotel-30.out;sentiment-30.out"
LATENCY_RESULTS_30="$DIR/latency-results-30.csv"
LATENCY_FILES_40="kn-full-election-40.csv kn-original-election-40.csv kn-partial-election-40.csv kn-full-hotel-40.csv kn-original-hotel-40.csv kn-partial-hotel-40.csv kn-full-sentiment-40.csv kn-original-sentiment-40.csv kn-partial-sentiment-40.csv"
LATENCY_OUT_40="election-40.out;hotel-40.out;sentiment-40.out"
LATENCY_RESULTS_40="$DIR/latency-results-40.csv"
LATENCY_FILES_50="kn-full-election-50.csv kn-original-election-50.csv kn-partial-election-50.csv kn-full-hotel-50.csv kn-original-hotel-50.csv kn-partial-hotel-50.csv kn-full-sentiment-50.csv kn-original-sentiment-50.csv kn-partial-sentiment-50.csv"
LATENCY_OUT_50="election-50.out;hotel-50.out;sentiment-50.out"
LATENCY_RESULTS_50="$DIR/latency-results-50.csv"
LATENCY_FILES_60="kn-full-election-60.csv kn-original-election-60.csv kn-partial-election-60.csv kn-full-hotel-60.csv kn-original-hotel-60.csv kn-partial-hotel-60.csv kn-full-sentiment-60.csv kn-original-sentiment-60.csv kn-partial-sentiment-60.csv"
LATENCY_OUT_60="election-60.out;hotel-60.out;sentiment-60.out"
LATENCY_RESULTS_60="$DIR/latency-results-60.csv"
LATENCY_FILES_70="kn-full-election-70.csv kn-original-election-70.csv kn-partial-election-70.csv kn-full-hotel-70.csv kn-original-hotel-70.csv kn-partial-hotel-70.csv kn-full-sentiment-70.csv kn-original-sentiment-70.csv kn-partial-sentiment-70.csv"
LATENCY_OUT_70="election-70.out;hotel-70.out;sentiment-70.out"
LATENCY_RESULTS_70="$DIR/latency-results-70.csv"
LATENCY_FILES_80="kn-full-election-80.csv kn-original-election-80.csv kn-partial-election-80.csv kn-full-hotel-80.csv kn-original-hotel-80.csv kn-partial-hotel-80.csv kn-full-sentiment-80.csv kn-original-sentiment-80.csv kn-partial-sentiment-80.csv"
LATENCY_OUT_80="election-80.out;hotel-80.out;sentiment-80.out"
LATENCY_RESULTS_80="$DIR/latency-results-80.csv"
LATENCY_FILES_90="kn-full-election-90.csv kn-original-election-90.csv kn-partial-election-90.csv kn-full-hotel-90.csv kn-original-hotel-90.csv kn-partial-hotel-90.csv kn-full-sentiment-90.csv kn-original-sentiment-90.csv kn-partial-sentiment-90.csv"
LATENCY_OUT_90="election-90.out;hotel-90.out;sentiment-90.out"
LATENCY_RESULTS_90="$DIR/latency-results-90.csv"
LATENCY_FILES_100="kn-full-election-100.csv kn-original-election-100.csv kn-partial-election-100.csv kn-full-hotel-100.csv kn-original-hotel-100.csv kn-partial-hotel-100.csv kn-full-sentiment-100.csv kn-original-sentiment-100.csv kn-partial-sentiment-100.csv"
LATENCY_OUT_100="election-100.out;hotel-100.out;sentiment-100.out"
LATENCY_RESULTS_100="$DIR/latency-results-100.csv"
LATENCY_FILES_110="kn-full-election-110.csv kn-original-election-110.csv kn-partial-election-110.csv kn-full-hotel-110.csv kn-original-hotel-110.csv kn-partial-hotel-110.csv kn-full-sentiment-110.csv kn-original-sentiment-110.csv kn-partial-sentiment-110.csv"
LATENCY_OUT_110="election-110.out;hotel-110.out;sentiment-110.out"
LATENCY_RESULTS_110="$DIR/latency-results-110.csv"
LATENCY_FILES_120="kn-full-election-120.csv kn-original-election-120.csv kn-partial-election-120.csv kn-full-hotel-120.csv kn-original-hotel-120.csv kn-partial-hotel-120.csv kn-full-sentiment-120.csv kn-original-sentiment-120.csv kn-partial-sentiment-120.csv"
LATENCY_OUT_120="election-120.out;hotel-120.out;sentiment-120.out"
LATENCY_RESULTS_120="$DIR/latency-results-120.csv"
LATENCY_FILES_130="kn-full-election-130.csv kn-original-election-130.csv kn-partial-election-130.csv kn-full-hotel-130.csv kn-original-hotel-130.csv kn-partial-hotel-130.csv kn-full-sentiment-130.csv kn-original-sentiment-130.csv kn-partial-sentiment-130.csv"
LATENCY_OUT_130="election-130.out;hotel-130.out;sentiment-130.out"
LATENCY_RESULTS_130="$DIR/latency-results-130.csv"
LATENCY_FILES_140="kn-full-election-140.csv kn-original-election-140.csv kn-partial-election-140.csv kn-full-hotel-140.csv kn-original-hotel-140.csv kn-partial-hotel-140.csv kn-full-sentiment-140.csv kn-original-sentiment-140.csv kn-partial-sentiment-140.csv"
LATENCY_OUT_140="election-140.out;hotel-140.out;sentiment-140.out"
LATENCY_RESULTS_140="$DIR/latency-results-140.csv"
LATENCY_FILES_150="kn-full-election-150.csv kn-original-election-150.csv kn-partial-election-150.csv kn-full-hotel-150.csv kn-original-hotel-150.csv kn-partial-hotel-150.csv kn-full-sentiment-150.csv kn-original-sentiment-150.csv kn-partial-sentiment-150.csv"
LATENCY_OUT_150="election-150.out;hotel-150.out;sentiment-150.out"
LATENCY_RESULTS_150="$DIR/latency-results-150.csv"
LATENCY_FILES_160="kn-full-election-160.csv kn-original-election-160.csv kn-partial-election-160.csv kn-full-hotel-160.csv kn-original-hotel-160.csv kn-partial-hotel-160.csv kn-full-sentiment-160.csv kn-original-sentiment-160.csv kn-partial-sentiment-160.csv"
LATENCY_OUT_160="election-160.out;hotel-160.out;sentiment-160.out"
LATENCY_RESULTS_160="$DIR/latency-results-160.csv"
LATENCY_FILES_170="kn-full-election-170.csv kn-original-election-170.csv kn-partial-election-170.csv kn-full-hotel-170.csv kn-original-hotel-170.csv kn-partial-hotel-170.csv kn-full-sentiment-170.csv kn-original-sentiment-170.csv kn-partial-sentiment-170.csv"
LATENCY_OUT_170="election-170.out;hotel-170.out;sentiment-170.out"
LATENCY_RESULTS_170="$DIR/latency-results-170.csv"
LATENCY_FILES_180="kn-full-election-180.csv kn-original-election-180.csv kn-partial-election-180.csv kn-full-hotel-180.csv kn-original-hotel-180.csv kn-partial-hotel-180.csv kn-full-sentiment-180.csv kn-original-sentiment-180.csv kn-partial-sentiment-180.csv"
LATENCY_OUT_180="election-180.out;hotel-180.out;sentiment-180.out"
LATENCY_RESULTS_180="$DIR/latency-results-180.csv"
LATENCY_FILES_190="kn-full-election-190.csv kn-original-election-190.csv kn-partial-election-190.csv kn-full-hotel-190.csv kn-original-hotel-190.csv kn-partial-hotel-190.csv kn-full-sentiment-190.csv kn-original-sentiment-190.csv kn-partial-sentiment-190.csv"
LATENCY_OUT_190="election-190.out;hotel-190.out;sentiment-190.out"
LATENCY_RESULTS_190="$DIR/latency-results-190.csv"
LATENCY_FILES_200="kn-full-election-200.csv kn-original-election-200.csv kn-partial-election-200.csv kn-full-hotel-200.csv kn-original-hotel-200.csv kn-partial-hotel-200.csv kn-full-sentiment-200.csv kn-original-sentiment-200.csv kn-partial-sentiment-200.csv"
LATENCY_OUT_200="election-200.out;hotel-200.out;sentiment-200.out"
LATENCY_RESULTS_200="$DIR/latency-results-200.csv"
LATENCY_FILES_210="kn-full-election-210.csv kn-original-election-210.csv kn-partial-election-210.csv kn-full-hotel-210.csv kn-original-hotel-210.csv kn-partial-hotel-210.csv kn-full-sentiment-210.csv kn-original-sentiment-210.csv kn-partial-sentiment-210.csv"
LATENCY_OUT_210="election-210.out;hotel-210.out;sentiment-210.out"
LATENCY_RESULTS_210="$DIR/latency-results-210.csv"
LATENCY_FILES_220="kn-full-election-220.csv kn-original-election-220.csv kn-partial-election-220.csv kn-full-hotel-220.csv kn-original-hotel-220.csv kn-partial-hotel-220.csv kn-full-sentiment-220.csv kn-original-sentiment-220.csv kn-partial-sentiment-220.csv"
LATENCY_OUT_220="election-220.out;hotel-220.out;sentiment-220.out"
LATENCY_RESULTS_220="$DIR/latency-results-220.csv"
LATENCY_FILES_230="kn-full-election-230.csv kn-original-election-230.csv kn-partial-election-230.csv kn-full-hotel-230.csv kn-original-hotel-230.csv kn-partial-hotel-230.csv kn-full-sentiment-230.csv kn-original-sentiment-230.csv kn-partial-sentiment-230.csv"
LATENCY_OUT_230="election-230.out;hotel-230.out;sentiment-230.out"
LATENCY_RESULTS_230="$DIR/latency-results-230.csv"
LATENCY_FILES_240="kn-full-election-240.csv kn-original-election-240.csv kn-partial-election-240.csv kn-full-hotel-240.csv kn-original-hotel-240.csv kn-partial-hotel-240.csv kn-full-sentiment-240.csv kn-original-sentiment-240.csv kn-partial-sentiment-240.csv"
LATENCY_OUT_240="election-240.out;hotel-240.out;sentiment-240.out"
LATENCY_RESULTS_240="$DIR/latency-results-240.csv"
LATENCY_FILES_250="kn-full-election-250.csv kn-original-election-250.csv kn-partial-election-250.csv kn-full-hotel-250.csv kn-original-hotel-250.csv kn-partial-hotel-250.csv kn-full-sentiment-250.csv kn-original-sentiment-250.csv kn-partial-sentiment-250.csv"
LATENCY_OUT_250="election-250.out;hotel-250.out;sentiment-250.out"
LATENCY_RESULTS_250="$DIR/latency-results-250.csv"
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_1 $DIR $LATENCY_FILES_1
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_2 $DIR $LATENCY_FILES_2
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_4 $DIR $LATENCY_FILES_4
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_8 $DIR $LATENCY_FILES_8
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_16 $DIR $LATENCY_FILES_16
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_20 $DIR $LATENCY_FILES_20
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_30 $DIR $LATENCY_FILES_30
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_40 $DIR $LATENCY_FILES_40
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_50 $DIR $LATENCY_FILES_50
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_60 $DIR $LATENCY_FILES_60
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_70 $DIR $LATENCY_FILES_70
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_80 $DIR $LATENCY_FILES_80
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_90 $DIR $LATENCY_FILES_90
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_100 $DIR $LATENCY_FILES_100
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_110 $DIR $LATENCY_FILES_110
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_120 $DIR $LATENCY_FILES_120
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_130 $DIR $LATENCY_FILES_130
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_140 $DIR $LATENCY_FILES_140
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_150 $DIR $LATENCY_FILES_150
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_160 $DIR $LATENCY_FILES_160
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_170 $DIR $LATENCY_FILES_170
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_180 $DIR $LATENCY_FILES_180
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_190 $DIR $LATENCY_FILES_190
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_200 $DIR $LATENCY_FILES_200
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_210 $DIR $LATENCY_FILES_210
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_220 $DIR $LATENCY_FILES_220
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_230 $DIR $LATENCY_FILES_230
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_240 $DIR $LATENCY_FILES_240
go run ../../tools/collection/processing-scripts/process-latency.go 500 $LATENCY_RESULTS_250 $DIR $LATENCY_FILES_250

go run ../../tools/collection/processing-scripts/process-metrics-all.go $METRIC_RESULTS $METRIC_FILES

go run ../../tools/collection/processing-scripts/process-multi-stats.go $STATS_RESULTS $METRIC_RESULTS $DIR "$LATENCY_RESULTS_1;$LATENCY_OUT_1" "$LATENCY_RESULTS_2;$LATENCY_OUT_2" "$LATENCY_RESULTS_4;$LATENCY_OUT_4" "$LATENCY_RESULTS_8;$LATENCY_OUT_8" "$LATENCY_RESULTS_16;$LATENCY_OUT_16" "$LATENCY_RESULTS_20;$LATENCY_OUT_20" "$LATENCY_RESULTS_30;$LATENCY_OUT_30" "$LATENCY_RESULTS_40;$LATENCY_OUT_40" "$LATENCY_RESULTS_50;$LATENCY_OUT_50" "$LATENCY_RESULTS_60;$LATENCY_OUT_60" "$LATENCY_RESULTS_70;$LATENCY_OUT_70" "$LATENCY_RESULTS_80;$LATENCY_OUT_80" "$LATENCY_RESULTS_90;$LATENCY_OUT_90" "$LATENCY_RESULTS_100;$LATENCY_OUT_100" "$LATENCY_RESULTS_110;$LATENCY_OUT_110" "$LATENCY_RESULTS_120;$LATENCY_OUT_120" "$LATENCY_RESULTS_130;$LATENCY_OUT_130" "$LATENCY_RESULTS_140;$LATENCY_OUT_140" "$LATENCY_RESULTS_150;$LATENCY_OUT_150" "$LATENCY_RESULTS_160;$LATENCY_OUT_160" "$LATENCY_RESULTS_170;$LATENCY_OUT_170" "$LATENCY_RESULTS_180;$LATENCY_OUT_180" "$LATENCY_RESULTS_190;$LATENCY_OUT_190" "$LATENCY_RESULTS_200;$LATENCY_OUT_200" "$LATENCY_RESULTS_210;$LATENCY_OUT_210" "$LATENCY_RESULTS_220;$LATENCY_OUT_220" "$LATENCY_RESULTS_230;$LATENCY_OUT_230" "$LATENCY_RESULTS_240;$LATENCY_OUT_240" "$LATENCY_RESULTS_250;$LATENCY_OUT_250"
