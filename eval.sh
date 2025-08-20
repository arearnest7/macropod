#!/bin/bash
HOST=${1:-127.0.0.1}
DIR=${2:-results}
mkdir $DIR
definitions=("election" "hotel" "sentiment" "video")
for definition in ${definitions[@]}; do
	cd kn-scripts
	sudo ./deploy-$definition-unified.sh $HOST
	sudo ./deploy-$definition-original.sh $HOST
	cd ..
        sleep 300s
	id="$(curl -X POST -d @eval-definitions/$definition.json http://10.43.190.1:9000/eval/start)"
	curl http://10.43.190.1:9000/eval/metrics/$id > $DIR/$definition-metrics.csv
	curl http://10.43.190.1:9000/eval/latency/$id > $DIR/$definition-latency.csv
	curl http://10.43.190.1:9000/eval/summary/$id > $DIR/$definition-summary.csv
	cd kn-scripts
	sudo ./remove-$definition-unified.sh
	sudo ./remove-$definition-original.sh
	cd ..
        kubectl delete ns knative-functions
        kubectl create ns knative-functions
	sleep 600s
done;
