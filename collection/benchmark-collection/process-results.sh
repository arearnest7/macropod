#!/bin/bash
N=${1:-10000}
LATENCY_RESULTS=${2:-latency-results.csv}
METRIC_RESULTS=${3:-metric-results.csv}
STATS_RESULTS=${4:-stats-results.csv}
DIR=${5:-results}
LOGS_DIR=${6:-results-logs}
LATENCY_FILES=${7:-"$DIR/macropod-election.csv $DIR/kn-full-election.csv $DIR/kn-original-election.csv $DIR/kn-partial-election.csv $DIR/macropod-feature.csv $DIR/kn-full-feature.csv $DIR/kn-original-feature.csv $DIR/kn-partial-feature.csv $DIR/macropod-hotel.csv $DIR/kn-full-hotel.csv $DIR/kn-original-hotel.csv $DIR/kn-partial-hotel.csv $DIR/macropod-pipelined.csv $DIR/kn-full-pipelined.csv $DIR/kn-original-pipelined.csv $DIR/kn-partial-pipelined.csv $DIR/macropod-sentiment.csv $DIR/kn-full-sentiment.csv $DIR/kn-original-sentiment.csv $DIR/kn-partial-sentiment.csv $DIR/macropod-video.csv $DIR/kn-full-video.csv $DIR/kn-original-video.csv $DIR/kn-partial-video.csv $DIR/macropod-wage.csv $DIR/kn-full-wage.csv $DIR/kn-original-wage.csv $DIR/kn-partial-wage.csv"}
LOG_FILES=${8:-"election-gateway election-get-results election-vote-enqueuer election-vote-processor feature-extractor feature-orchestrator feature-reducer feature-status feature-wait hotel-frontend hotel-geo hotel-profile hotel-rate hotel-recommend hotel-reserve hotel-search hotel-user pipelined-checksum pipelined-encrypt pipelined-main pipelined-zip sentiment-cfail sentiment-db sentiment-main sentiment-product-or-service sentiment-product-result sentiment-product-sentiment sentiment-read-csv sentiment-service-result sentiment-service-sentiment sentiment-sfail sentiment-sns video-decoder video-recog video-streaming wage-avg wage-format wage-merit wage-stats wage-sum wage-validator wage-write-merit wage-write-raw election-full feature-full hotel-full pipelined-full sentiment-full video-full wage-full"}
METRIC_FILES=${9:-"$DIR/node1-metrics.csv $DIR/node2-metrics.csv $DIR/node3-metrics.csv $DIR/node4-metrics.csv $DIR/node5-metrics.csv"}
MACROPOD_LOG_BUNDLES=${10:-"macropod-election;macropod-election-log-bundle.csv;election-gateway;election-get-results;election-vote-enqueuer;election-vote-processor macropod-feature;macropod-feature-log-bundle.csv;feature-extractor;feature-orchestrator;feature-reducer;feature-status;feature-wait macropod-hotel;macropod-hotel-log-bundle.csv;hotel-frontend;hotel-geo;hotel-rate;hotel-recommend;hotel-reserve;hotel-search;hotel-user macropod-pipelined;macropod-pipelined-log-bundle.csv;pipelined-checksum;pipelined-encrypt;pipelined-main;pipelined-zip macropod-sentiment;macropod-sentiment-log-bundle.csv;sentiment-cfail;sentiment-db;sentiment-main;sentiment-product-or-service;sentiment-product-result;sentiment-product-sentiment;sentiment-read-csv;sentiment-service-result;sentiment-service-sentiment;sentiment-sfail;sentiment-sns macropod-video;macropod-video-log-bundle.csv;video-decoder;video-recog;video-streaming macropod-wage;macropod-wage-log-bundle.csv;wage-avg;wage-format;wage-merit;wage-stats;wage-sum;wage-validator;wage-write-merit;wage-write-raw"}
mkdir $LOGS_DIR
go run ../../tools/processing-scripts/process-latency.go $N $LATENCY_RESULTS $LATENCY_FILES
go run ../../tools/processing-scripts/split-macropod-log-bundles.go $DIR $MACROPOD_LOG_BUNDLES
go run ../../tools/processing-scripts/process-metrics.go $METRIC_RESULTS $METRIC_FILES
go run ../../tools/processing-scripts/clean-logs.go $DIR $LOGS_DIR $LOG_FILES

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-election-processed-log.csv kn-full-election-election full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-election-processed-log.csv kn-original-election-election gateway vote-enqueuer vote-processor
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-election-processed-log.csv kn-partial-election-election gateway-vevp
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-election-processed-log.csv macropod-election-election gateway vote-enqueuer vote-processor

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-feature-processed-log.csv kn-full-feature-feature full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-feature-processed-log.csv kn-original-feature-feature orchestrator extractor wait status reducer
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-feature-processed-log.csv kn-partial-feature-feature orchestrator-wsr extractor-partial
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-feature-processed-log.csv macropod-feature-feature orchestrator extractor wait status reducer

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-hotel-processed-log.csv kn-full-hotel-hotel full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-hotel-processed-log.csv kn-original-hotel-hotel frontend search geo rate
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-hotel-processed-log.csv kn-partial-hotel-hotel frontend-spgr
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-hotel-processed-log.csv macropod-hotel-hotel frontend search geo rate

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-pipelined-processed-log.csv kn-full-pipelined-pipelined full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-pipelined-processed-log.csv kn-original-pipelined-pipelined main checksum zip encrypt
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-pipelined-processed-log.csv kn-partial-pipelined-pipelined main-partial checksum-partial zip-partial encrypt-partial
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-pipelined-processed-log.csv macropod-pipelined-pipelined main checksum zip encrypt

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-sentiment-processed-log.csv kn-full-sentiment-sentiment full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-sentiment-processed-log.csv kn-original-sentiment-sentiment main read-csv product-or-service product-sentiment product-result sns db
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-sentiment-processed-log.csv kn-partial-sentiment-sentiment main-rcposc product-sentiment-prs db-s
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-sentiment-processed-log.csv macropod-sentiment-sentiment main read-csv product-or-service product-sentiment product-result sns db

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-video-processed-log.csv kn-full-video-video full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-video-processed-log.csv kn-original-video-video streaming decoder recog
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-video-processed-log.csv kn-partial-video-video streaming-d recog-partial
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-video-processed-log.csv streaming decoder recog

go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-full-wage-processed-log.csv kn-full-wage-wage full
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-original-wage-processed-log.csv kn-original-wage-wage validator format write-raw stats sum avg merit write-merit
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR kn-partial-wage-processed-log.csv kn-partial-wage-wage validator-fw stats-partial sum-amw
go run ../../tools/processing-scripts/process-logs.go "0:0:0" $LOGS_DIR macropod-wage-processed-log.csv validator format write-raw stats sum avg merit write-merit

go run ../../tools/processing-scripts/process-stats.go $STATS_RESULTS $LATENCY_RESULTS $METRIC_RESULTS "macropod-election;macropod-election-processed-log.csv" "kn-full-election;kn-full-election-processed-log.csv" "kn-original-election;kn-original-election-processed-log.csv" "kn-partial-election;kn-partial-election-processed-log.csv" "macropod-feature;macropod-feature-processed-log.csv" "kn-full-feature;kn-full-feature-processed-log.csv" "kn-original-feature;kn-original-feature-processed-log.csv" "kn-partial-feature;kn-partial-feature-processed-log.csv" "macropod-pod-hotel;macropod-hotel-processed-log.csv" "kn-full-hotel;kn-full-hotel-processed-log.csv" "kn-original-hotel;kn-original-hotel-processed-log.csv" "kn-partial-hotel;kn-partial-hotel-processed-log.csv" "macropod-pipelined;macropod-pipelined-processed-log.csv" "kn-full-pipelined;kn-full-pipelined-processed-log.csv" "kn-original-pipelined;kn-original-pipelined-processed-log.csv" "kn-partial-pipelined;kn-partial-pipelined-processed-log.csv" "macropod-sentiment;macropod-sentiment-processed-log.csv" "kn-full-sentiment;kn-full-sentiment-processed-log.csv" "kn-original-sentiment;kn-original-sentiment-processed-log.csv" "kn-partial-sentiment;kn-partial-sentiment-processed-log.csv" "macropod-video;macropod-video-processed-log.csv" "kn-full-video;kn-full-video-processed-log.csv" "kn-original-video;kn-original-video-processed-log.csv" "kn-partial-video;kn-partial-video-processed-log.csv" "macropod-wage;macropod-wage-processed-log.csv" "kn-full-wage;kn-full-wage-processed-log.csv" "kn-original-wage;kn-original-wage-processed-log.csv" "kn-partial-wage;kn-partial-wage-processed-log.csv"
