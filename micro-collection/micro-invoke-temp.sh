#!/bin/bash
DEST=${1:-1MB}
SIZE=${2:-1000000}
sed -i "s/1000000/$SIZE/g" micro-rpc-grpc.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-pod.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-faastroute.yaml
kubectl apply -f micro-rpc-grpc.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv http://10.43.190.1 >> micro-rpc-grpc-$DEST.csv
kubectl delete -f micro-rpc-grpc.yaml
sleep 1000s
kubectl apply -f micro-rpc-pod.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv http://10.43.190.1 >> micro-rpc-pod-$DEST.csv
kubectl delete -f micro-rpc-pod.yaml
sleep 1000s
kubectl apply -f micro-rpc-faastroute.yaml
sleep 1000s
hey -n 10000 -c 1 -t 30 -o csv http://10.43.190.1 >> micro-rpc-faastroute-$DEST.csv
kubectl delete -f micro-rpc-faastroute.yaml
sleep 1000s
sed -i "s/$SIZE/1000000/g" micro-rpc-grpc.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-pod.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-faastroute.yaml
