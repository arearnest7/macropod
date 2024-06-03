#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-../payloads/election.json}
N=${4:-10000}
containers=(${5:-"election-gateway election-get-results election-vote-enqueuer election-vote-processor"})
cd ../../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-multi-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/benchmark-collection/multi-oci-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> ../../collection/benchmark-collection/multi-oci-$i.csv; done;
kubectl delete -f ./$BENCHMARK-multi-oci.yaml
sleep 180s
kubectl apply -f ./$BENCHMARK-multi-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/benchmark-collection/multi-pod-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> ../../collection/benchmark-collection/multi-pod-$i.csv; done;
kubectl delete -f ./$BENCHMARK-multi-pod.yaml
sleep 180s
kubectl apply -f ./$BENCHMARK-single-mmap.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/benchmark-collection/single-mmap-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do for container in ${containers[@]}; do kubectl logs $i -c $container >> ../../collection/benchmark-collection/single-mmap-$container-$i.csv; done; done;
kubectl delete -f ./$BENCHMARK-single-mmap.yaml
sleep 180s
kubectl apply -f ./$BENCHMARK-single-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/benchmark-collection/single-oci-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do for container in ${containers[@]}; do kubectl logs $i -c $container >> ../../collection/benchmark-collection/single-oci-$container-$i.csv; done; done;
kubectl delete -f ./$BENCHMARK-single-oci.yaml
sleep 180s
kubectl apply -f ./$BENCHMARK-single-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/benchmark-collection/single-pod-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do for container in ${containers[@]}; do kubectl logs $i -c $container >> ../../collection/benchmark-collection/single-pod-$container-$i.csv; done; done;
kubectl delete -f ./$BENCHMARK-single-pod.yaml
cd ../../collection/benchmark-collection/
sleep 180s
