#!/bin/bash
HOST=${1:-10.43.190.1}
DIR=${2:-results}
mkdir $DIR
definitions=("excamera" "finra" "hotel-reservation" "sentiment-analysis" "serverless-election" "social-network" "video-analytics" "wage-pay")
for definition in ${definitions[@]}; do
	id="$(curl -X POST -d @eval-definitions/$definition.json http://$HOST:9000/eval/start)"
	curl http://$HOST:9000/eval/metrics/$id > $DIR/$definition-metrics.csv
	curl http://$HOST:9000/eval/latency/$id > $DIR/$definition-latency.csv
	curl http://$HOST:9000/eval/summary/$id > $DIR/$definition-summary.csv
	sleep 60s
done;
