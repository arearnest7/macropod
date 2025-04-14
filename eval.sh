#!/bin/bash
target=${1:-10.43.190.1}
results=${2:-results}
agg=(election-agg hotel-agg pipelined-agg sentiment-agg video-agg wage-agg)
disagg-w-gateway=(election-disagg-w-gateway hotel-disagg-w-gateway pipelined-disagg-w-gateway sentiment-disagg-w-gateway video-disagg-w-gateway wage-disagg-w-gateway)
disagg-wo-gateway=(election-disagg-wo-gateway hotel-disagg-wo-gateway pipelined-disagg-wo-gateway sentiment-disagg-wo-gateway video-disagg-wo-gateway wage-disagg-wo-gateway)
dynamic=(election-dynamic hotel-dynamic pipelined-dynamic sentiment-dynamic video-dynamic wage-dynamic)
unified=(election-unified hotel-unified pipelined-unified sentiment-unified video-unified wage-unified)

mkdir $results
for i in "${agg[@]}"; do ./eval-invoke.sh $target $i $results; done;
for i in "${disagg-w-gateway[@]}"; do ./eval-invoke.sh $target $i $results; done;
for i in "${disagg-wo-gateway[@]}"; do ./eval-invoke.sh $target $i $results; done;
for i in "${dynamic[@]}"; do ./eval-invoke.sh $target $i $results; done;
for i in "${unified[@]}"; do ./eval-invoke.sh $target $i $results; done;
