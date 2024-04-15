#!/bin/bash
DEST=${1:-1MB}
SIZE=${2:-1000000}
HOST=${3:-localhost}
LOGGING_HOST=${4:-localhost}
sed -i "s/1000000/$SIZE/g" ../benchmarks/micro/micro-rpc-a/func.yaml
sed -i "s/1000000/$SIZE/g" ../benchmarks/micro/micro-rpc-b/func.yaml
sed -i "s/1000000/$SIZE/g" ../benchmarks/micro/micro-rpc-a-b/func.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-multi-pod.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-single-pod.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-single-mmap.yaml
kn func deploy --build=false --push=false --path ../benchmarks/micro/micro-rpc-a
kn func deploy --build=false --push=false --path ../benchmarks/micro/micro-rpc-b
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv -m POST -T application/json http://micro-rpc-a.default.$HOST.sslip.io >> micro-rpc-kn-original-$DEST.csv
kn func delete micro-rpc-a
kn func delete micro-rpc-b
sleep 1000s
kn func deploy --build=false --push=false --path ../benchmarks/micro/micro-rpc-a-b
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv -m POST -T application/json http://micro-rpc-a-b.default.$HOST.sslip.io >> micro-rpc-kn-full-$DEST.csv
kn func delete micro-rpc-a-b
sleep 1000s
kubectl apply -f micro-rpc-multi-pod.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-multi-pod-$DEST.csv
kubectl delete -f micro-rpc-multi-pod.yaml
sleep 1000s
kubectl apply -f micro-rpc-single-pod.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-single-pod-$DEST.csv
kubectl delete -f micro-rpc-single-pod.yaml
sleep 1000s
kubectl apply -f micro-rpc-single-mmap.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-single-mmap-$DEST.csv
kubectl delete -f micro-rpc-single-mmap.yaml
sleep 1000s
sed -i "s/$SIZE/1000000/g" ../benchmarks/micro/micro-rpc-a/func.yaml
sed -i "s/$SIZE/1000000/g" ../benchmarks/micro/micro-rpc-b/func.yaml
sed -i "s/$SIZE/1000000/g" ../benchmarks/micro/micro-rpc-a-b/func.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-multi-pod.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-single-pod.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-single-mmap.yaml
