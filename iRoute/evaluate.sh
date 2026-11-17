#!/bin/bash
HOST=${1:-10.43.190.1}
IROUTE=${2:-192.168.0.1}
DIR=${3:-results}
mkdir $DIR
definitions=("excamera" "finra" "social-network")
for definition in ${definitions[@]}; do
	cd scripts/deploy
	sudo ./deploy-$definition.sh
	cd ../..
	sed -i "s/192\.168\.0\.1/$IROUTE/g" eval-definitions/$definition.json
        sleep 60s
	id="$(curl -X POST -d @eval-definitions/$definition.json http://$HOST:9000/eval/start)"
	curl http://$HOST:9000/eval/metrics/$id > $DIR/$definition-metrics.csv
	curl http://$HOST:9000/eval/latency/$id > $DIR/$definition-latency.csv
	curl http://$HOST:9000/eval/summary/$id > $DIR/$definition-summary.csv
	sed -i "s/$IROUTE/192\.168\.0\.1/g" eval-definitions/$definition.json
	cd scripts/remove
	sudo ./remove-$definition.sh
	cd ../..
done;
