#!/bin/bash
BENCHMARK=${1:-election}
C=${2:-1}
PAYLOAD=${3:-../payloads/election.json}
N=${4:-10000}
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-multi-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> multi-oci-$BENCHMARK.csv
mv multi-oci-$BENCHMARK.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> multi-oci-$i.csv; done;
kubectl delete -f ./$BENCHMARK-multi-oci.yaml
cd ../../collection
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-multi-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> multi-pod-$BENCHMARK.csv
mv multi-pod-$BENCHMARK.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> multi-pod-$i.csv; done;
kubectl delete -f ./$BENCHMARK-multi-pod.yaml
cd ../../collection
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-mmap.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-mmap-$BENCHMARK.csv
mv single-mmap-$BENCHMARK.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> single-mmap-$i.csv; done;
kubectl delete -f ./$BENCHMARK-single-mmap.yaml
cd ../../collection
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-oci-$BENCHMARK.csv
mv single-oci-$BENCHMARK.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> single-oci-$i.csv; done;
kubectl delete -f ./$BENCHMARK-single-oci.yaml
cd ../../collection
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-pod-$BENCHMARK.csv
mv single-pod-$BENCHMARK.csv ../../collection/
logs=$(kubectl get pods --no-headers -o custom-columns=":metadata.name" --sort-by="metadata.name")
for i in $logs; do kubectl logs $i >> single-pod-$i.csv; done;
kubectl delete -f ./$BENCHMARK-single-pod.yaml
cd ../../collection
sleep 180s
