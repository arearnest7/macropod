#!/bin/bash
N=${1:-10000}
LATENCY_RESULTS=${2:-latency-results.csv}
METRIC_RESULTS=${3:-metric-results.csv}
STATS_RESULTS=${4:-stats-results.csv}
DIR=${5:-results}
LOGS_DIR=${6:-results-logs}
LATENCY_FILES=${7:-"$DIR/multi-oci-election.csv $DIR/multi-pod-election.csv $DIR/single-mmap-election.csv $DIR/single-oci-election.csv $DIR/single-pod-election.csv $DIR/kn-full-election.csv $DIR/kn-original-election.csv $DIR/kn-partial-election.csv $DIR/multi-oci-feature.csv $DIR/multi-pod-feature.csv $DIR/single-mmap-feature.csv $DIR/single-oci-feature.csv $DIR/single-pod-feature.csv $DIR/kn-full-feature.csv $DIR/kn-original-feature.csv $DIR/kn-partial-feature.csv $DIR/multi-oci-hotel.csv $DIR/multi-pod-hotel.csv $DIR/single-mmap-hotel.csv $DIR/single-oci-hotel.csv $DIR/single-pod-hotel.csv $DIR/kn-full-hotel.csv $DIR/kn-original-hotel.csv $DIR/kn-partial-hotel.csv $DIR/multi-oci-pipelined.csv $DIR/multi-pod-pipelined.csv $DIR/single-mmap-pipelined.csv $DIR/single-oci-pipelined.csv $DIR/single-pod-pipelined.csv $DIR/kn-full-pipelined.csv $DIR/kn-original-pipelined.csv $DIR/kn-partial-pipelined.csv $DIR/multi-oci-sentiment.csv $DIR/multi-pod-sentiment.csv $DIR/single-mmap-sentiment.csv $DIR/single-oci-sentiment.csv $DIR/single-pod-sentiment.csv $DIR/kn-full-sentiment.csv $DIR/kn-original-sentiment.csv $DIR/kn-partial-sentiment.csv $DIR/multi-oci-video.csv $DIR/multi-pod-video.csv $DIR/single-mmap-video.csv $DIR/single-oci-video.csv $DIR/single-pod-video.csv $DIR/kn-full-video.csv $DIR/kn-original-video.csv $DIR/kn-partial-video.csv $DIR/multi-oci-wage.csv $DIR/multi-pod-wage.csv $DIR/single-mmap-wage.csv $DIR/single-oci-wage.csv $DIR/single-pod-wage.csv $DIR/kn-full-wage.csv $DIR/kn-original-wage.csv $DIR/kn-partial-wage.csv"}
LOG_FILES=${8:-"election-gateway election-get-results election-vote-enqueuer election-vote-processor feature-extractor feature-orchestrator feature-reducer feature-status feature-wait hotel-frontend hotel-geo hotel-profile hotel-rate hotel-recommend hotel-reserve hotel-search hotel-user pipelined-checksum pipelined-encrypt pipelined-main pipelined-zip sentiment-cfail sentiment-db sentiment-main sentiment-product-or-service sentiment-product-result sentiment-product-sentiment sentiment-read-csv sentiment-service-result sentiment-service-sentiment sentiment-sfail sentiment-sns video-decoder video-recog video-streaming wage-avg wage-format wage-merit wage-stats wage-sum wage-validator wage-write-merit wage-write-raw election-full feature-full hotel-full pipelined-full sentiment-full video-full wage-full"}
METRIC_FILES=${9:-"$DIR/node1-metrics.csv $DIR/node2-metrics.csv $DIR/node3-metrics.csv $DIR/node4-metrics.csv $DIR/node5-metrics.csv"}
mkdir $LOGS_DIR
go run ../../tools/processing-scripts/process-latency.go $N $LATENCY_RESULTS $LATENCY_FILES
go run ../../tools/processing-scripts/clean-logs.go $DIR $LOGS_DIR $LOG_FILES
go run ../../tools/processing-scripts/process-metrics.go $METRIC_RESULTS $METRIC_FILES

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-election-processed-log.csv multi-oci-election gateway vote-enqueuer vote-processor
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-election-processed-log.csv multi-pod-election gateway vote-enqueuer vote-processor
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-election-processed-log.csv single-mmap-election election-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-election-processed-log.csv single-oci-election election-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-election-processed-log.csv single-pod-election election-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-election-processed-log.csv kn-full-election-election full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-election-processed-log.csv kn-original-election-election gateway vote-enqueuer vote-processor
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-election-processed-log.csv kn-partial-election-election gateway-vevp

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-feature-processed-log.csv multi-oci-feature orchestrator extractor wait status reducer
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-feature-processed-log.csv multi-pod-feature orchestrator extractor wait status reducer
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-feature-processed-log.csv single-mmap-feature feature-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-feature-processed-log.csv single-oci-feature feature-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-feature-processed-log.csv single-pod-feature feature-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-feature-processed-log.csv kn-full-feature-feature full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-feature-processed-log.csv kn-original-feature-feature orchestrator extractor wait status reducer
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-feature-processed-log.csv kn-partial-feature-feature orchestrator-wsr extractor-partial

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-hotel-processed-log.csv multi-oci-hotel frontend search geo rate
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-hotel-processed-log.csv multi-pod-hotel frontend search geo rate
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-hotel-processed-log.csv single-mmap-hotel hotel-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-hotel-processed-log.csv single-oci-hotel hotel-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-hotel-processed-log.csv single-pod-hotel hotel-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-hotel-processed-log.csv kn-full-hotel-hotel full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-hotel-processed-log.csv kn-original-hotel-hotel frontend search geo rate
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-hotel-processed-log.csv kn-partial-hotel-hotel frontend-spgr

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-pipelined-processed-log.csv multi-oci-pipelined main checksum zip encrypt
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-pipelined-processed-log.csv multi-pod-pipelined main checksum zip encrypt
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-pipelined-processed-log.csv single-mmap-pipelined pipelined-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-pipelined-processed-log.csv single-oci-pipelined pipelined-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-pipelined-processed-log.csv single-pod-pipelined pipelined-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-pipelined-processed-log.csv kn-full-pipelined-pipelined full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-pipelined-processed-log.csv kn-original-pipelined-pipelined main checksum zip encrypt
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-pipelined-processed-log.csv kn-partial-pipelined-pipelined main-partial checksum-partial zip-partial encrypt-partial

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-sentiment-processed-log.csv multi-oci-sentiment main read-csv product-or-service product-sentiment product-result sns db
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-sentiment-processed-log.csv multi-pod-sentiment main read-csv product-or-service product-sentiment product-result sns db
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-sentiment-processed-log.csv single-mmap-sentiment sentiment-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-sentiment-processed-log.csv single-oci-sentiment sentiment-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-sentiment-processed-log.csv single-pod-sentiment sentiment-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-sentiment-processed-log.csv kn-full-sentiment-sentiment full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-sentiment-processed-log.csv kn-original-sentiment-sentiment main read-csv product-or-service product-sentiment product-result sns db
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-sentiment-processed-log.csv kn-partial-sentiment-sentiment main-rcposc product-sentiment-prs db-s

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-video-processed-log.csv multi-oci-video streaming decoder recog
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-video-processed-log.csv multi-pod-video streaming decoder recog
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-video-processed-log.csv single-mmap-video video-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-video-processed-log.csv single-oci-video video-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-video-processed-log.csv single-pod-video video-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-video-processed-log.csv kn-full-video-video full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-video-processed-log.csv kn-original-video-video streaming decoder recog
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-video-processed-log.csv kn-partial-video-video streaming-d recog-partial

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-oci-wage-processed-log.csv multi-oci-wage validator format write-raw stats sum avg merit write-merit
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR multi-pod-wage-processed-log.csv multi-pod-wage validator format write-raw stats sum avg merit write-merit
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-mmap-wage-processed-log.csv single-mmap-wage wage-single-mmap
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-oci-wage-processed-log.csv single-oci-wage wage-single-oci
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR single-pod-wage-processed-log.csv single-pod-wage wage-single-pod
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-wage-processed-log.csv kn-full-wage-wage full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-wage-processed-log.csv kn-original-wage-wage validator format write-raw stats sum avg merit write-merit
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-wage-processed-log.csv kn-partial-wage-wage validator-fw stats-partial sum-amw

go run ../../tools/processing-scripts/process-stats.go $STATS_RESULTS $LATENCY_RESULTS $METRIC_RESULTS "multi-oci-election;multi-oci-election-processed-log.csv" "multi-pod-election;multi-pod-election-processed-log.csv" "single-mmap-election;single-mmap-election-processed-log.csv" "single-oci-election;single-oci-election-processed-log.csv" "single-pod-election;single-pod-election-processed-log.csv" "kn-full-election;kn-full-election-processed-log.csv" "kn-original-election;kn-original-election-processed-log.csv" "kn-partial-election;kn-partial-election-processed-log.csv" "multi-oci-feature;multi-oci-feature-processed-log.csv" "multi-pod-feature;multi-pod-feature-processed-log.csv" "single-mmap-feature;single-mmap-feature-processed-log.csv" "single-oci-feature;single-oci-feature-processed-log.csv" "single-pod-feature;single-pod-feature-processed-log.csv" "kn-full-feature;kn-full-feature-processed-log.csv" "kn-original-feature;kn-original-feature-processed-log.csv" "kn-partial-feature;kn-partial-feature-processed-log.csv" "multi-oci-hotel;multi-oci-hotel-processed-log.csv" "multi-pod-hotel;multi-pod-hotel-processed-log.csv" "single-mmap-hotel;single-mmap-hotel-processed-log.csv" "single-oci-hotel;single-oci-hotel-processed-log.csv" "single-pod-hotel;single-pod-hotel-processed-log.csv" "kn-full-hotel;kn-full-hotel-processed-log.csv" "kn-original-hotel;kn-original-hotel-processed-log.csv" "kn-partial-hotel;kn-partial-hotel-processed-log.csv" "multi-oci-pipelined;multi-oci-pipelined-processed-log.csv" "multi-pod-pipelined;multi-pod-pipelined-processed-log.csv" "single-mmap-pipelined;single-mmap-pipelined-processed-log.csv" "single-oci-pipelined;single-oci-pipelined-processed-log.csv" "single-pod-pipelined;single-pod-pipelined-processed-log.csv" "kn-full-pipelined;kn-full-pipelined-processed-log.csv" "kn-original-pipelined;kn-original-pipelined-processed-log.csv" "kn-partial-pipelined;kn-partial-pipelined-processed-log.csv" "multi-oci-sentiment;multi-oci-sentiment-processed-log.csv" "multi-pod-sentiment;multi-pod-sentiment-processed-log.csv" "single-mmap-sentiment;single-mmap-sentiment-processed-log.csv" "single-oci-sentiment;single-oci-sentiment-processed-log.csv" "single-pod-sentiment;single-pod-sentiment-processed-log.csv" "kn-full-sentiment;kn-full-sentiment-processed-log.csv" "kn-original-sentiment;kn-original-sentiment-processed-log.csv" "kn-partial-sentiment;kn-partial-sentiment-processed-log.csv" "multi-oci-video;multi-oci-video-processed-log.csv" "multi-pod-video;multi-pod-video-processed-log.csv" "single-mmap-video;single-mmap-video-processed-log.csv" "single-oci-video;single-oci-video-processed-log.csv" "single-pod-video;single-pod-video-processed-log.csv" "kn-full-video;kn-full-video-processed-log.csv" "kn-original-video;kn-original-video-processed-log.csv" "kn-partial-video;kn-partial-video-processed-log.csv" "multi-oci-wage;multi-oci-wage-processed-log.csv" "multi-pod-wage;multi-pod-wage-processed-log.csv" "single-mmap-wage;single-mmap-wage-processed-log.csv" "single-oci-wage;single-oci-wage-processed-log.csv" "single-pod-wage;single-pod-wage-processed-log.csv" "kn-full-wage;kn-full-wage-processed-log.csv" "kn-original-wage;kn-original-wage-processed-log.csv" "kn-partial-wage;kn-partial-wage-processed-log.csv"
