#!/bin/bash
BENCHMARK=${1:-election}
REDIS=${2:-127.0.0.1}
PASSWORD=${3:-password}
C=${4:-1}
PAYLOAD=${5:-../payloads/election.json}
N=${6:-10000}
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-multi-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> multi-oci-$BENCHMARK.csv
mv multi-oci-$BENCHMARK.csv ../../collection/
kubectl delete -f ./$BENCHMARK-multi-oci.yaml
cd ../../collection
./collect-redis-original-$BENCHMARK.sh $REDIS $PASSWORD multi-oci
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-multi-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> multi-pod-$BENCHMARK.csv
mv multi-pod-$BENCHMARK.csv ../../collection/
kubectl delete -f ./$BENCHMARK-multi-pod.yaml
cd ../../collection
./collect-redis-original-$BENCHMARK.sh $REDIS $PASSWORD multi-pod
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-mmap.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-mmap-$BENCHMARK.csv
mv single-mmap-$BENCHMARK.csv ../../collection/
kubectl delete -f ./$BENCHMARK-single-mmap.yaml
cd ../../collection
./collect-redis-original-$BENCHMARK.sh $REDIS $PASSWORD single-mmap
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-oci.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-oci-$BENCHMARK.csv
mv single-oci-$BENCHMARK.csv ../../collection/
kubectl delete -f ./$BENCHMARK-single-oci.yaml
cd ../../collection
./collect-redis-original-$BENCHMARK.sh $REDIS $PASSWORD single-oci
sleep 180s
cd ../tools/deployment-yamls
kubectl apply -f ./$BENCHMARK-single-pod.yaml
sleep 180s
hey -n $N -c $C -t 180 -o csv -D $PAYLOAD http://10.43.190.1 >> single-pod-$BENCHMARK.csv
mv single-pod-$BENCHMARK.csv ../../collection/
kubectl delete -f ./$BENCHMARK-single-pod.yaml
cd ../../collection
./collect-redis-original-$BENCHMARK.sh $REDIS $PASSWORD single-pod
sleep 180s
