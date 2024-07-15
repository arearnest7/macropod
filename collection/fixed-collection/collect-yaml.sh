#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-../payloads/election.json}
N=${4:-10000}
cd ../../tools/workflow
kubectl apply -f ./fixed-oci/$BENCHMARK-multi-oci.yaml
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/fixed-collection/multi-oci-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i --all-containers=true >> ../../collection/fixed-collection/multi-oci-$i.csv; done;
kubectl delete -f ./fixed-oci/$BENCHMARK-multi-oci.yaml
sleep 180s
kubectl apply -f ./fixed-macropod/$BENCHMARK-multi-pod.yaml
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/fixed-collection/multi-pod-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i --all-containers=true >> ../../collection/fixed-collection/multi-pod-$i.csv; done;
kubectl delete -f ./fixed-macropod/$BENCHMARK-multi-pod.yaml
sleep 180s
kubectl apply -f ./fixed-macropod/$BENCHMARK-single-mmap.yaml
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/fixed-collection/single-mmap-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i --all-containers=true >> ../../collection/fixed-collection/single-mmap-$i.csv; done;
kubectl delete -f ./fixed-macropod/$BENCHMARK-single-mmap.yaml
sleep 180s
kubectl apply -f ./fixed-oci/$BENCHMARK-single-oci.yaml
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/fixed-collection/single-oci-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i --all-containers=true >> ../../collection/fixed-collection/single-oci-$i.csv; done;
kubectl delete -f ./fixed-oci/$BENCHMARK-single-oci.yaml
sleep 180s
kubectl apply -f ./fixed-macropod/$BENCHMARK-single-pod.yaml
sleep 180s
hey -n $N -c $C -t 1000 -o csv -D $PAYLOAD -m POST -T application/json http://10.43.190.1 >> ../../collection/fixed-collection/single-pod-$BENCHMARK.csv
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i --all-containers=true >> ../../collection/fixed-collection/single-pod-$i.csv; done;
kubectl delete -f ./fixed-macropod/$BENCHMARK-single-pod.yaml
cd ../../collection/fixed-collection/
sleep 180s
