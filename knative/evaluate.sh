#!/bin/bash
HOST=${1:-10.43.190.1}
KNATIVE=${2:-192.168.0.1}
DIR=${3:-results}
mkdir $DIR
definitions=("excamera" "finra" "hotel-reservation" "sentiment-analysis" "serverless-election" "social-network" "video-analytics" "wage-pay")
for definition in ${definitions[@]}; do
	cd scripts/deploy
	sudo ./deploy-$definition.sh $KNATIVE
	cd ../..
	sed -i "s/192\.168\.0\.1/$KNATIVE/g" eval-definitions/$definition.json
        sleep 60s
	id="$(curl -X POST -d @eval-definitions/$definition.json http://$HOST:9000/eval/start)"
	curl http://$HOST:9000/eval/metrics/$id > $DIR/$definition-metrics.csv
	curl http://$HOST:9000/eval/latency/$id > $DIR/$definition-latency.csv
	curl http://$HOST:9000/eval/summary/$id > $DIR/$definition-summary.csv
        sed -i "s/$KNATIVE/192\.168\.0\.1/g" eval-definitions/$definition.json
	cd scripts/remove
	sudo ./remove-$definition.sh
	cd ../..
        kubectl delete ns knative-functions
        kubectl create ns knative-functions
done;
