#!/bin/bash
DEST=${1:-1MB}
SIZE=${2:-1000000}
HOST=${3:-127.0.0.1}
N=${4:-10000}
C=${5:-1}
sed -i "s/1000000/$SIZE/g" ../../benchmarks/micro/micro-rpc-a/func.yaml
sed -i "s/1000000/$SIZE/g" ../../benchmarks/micro/micro-rpc-b/func.yaml
sed -i "s/1000000/$SIZE/g" ../../benchmarks/micro/micro-rpc-a-b/func.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-multi-pod.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-single-pod.yaml
sed -i "s/1000000/$SIZE/g" micro-rpc-single-mmap.yaml
sudo kn func deploy --build=false --push=false --path ../../benchmarks/micro/micro-rpc-a
sudo kn func deploy --build=false --push=false --path ../../benchmarks/micro/micro-rpc-b
sleep 180s
hey -n $N -c $C -t 1000 -d "{}" -o csv -m POST -T application/json http://micro-rpc-a.knative-functions.$HOST.sslip.io >> micro-rpc-kn-original-$DEST.csv
sudo kn func delete micro-rpc-a
sudo kn func delete micro-rpc-b
sleep 180s
sudo kn func deploy --build=false --push=false --path ../../benchmarks/micro/micro-rpc-a-b
sleep 180s
hey -n $N -c $C -t 1000 -d "{}" -o csv -m POST -T application/json http://micro-rpc-a-b.knative-functions.$HOST.sslip.io >> micro-rpc-kn-full-$DEST.csv
sudo kn func delete micro-rpc-a-b
sleep 180s
sudo kubectl apply -f micro-rpc-multi-pod.yaml
sleep 180s
hey -n $N -c $C -t 1000 -d "{}" -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-multi-pod-$DEST.csv
sudo kubectl delete -f micro-rpc-multi-pod.yaml
sleep 180s
sudo kubectl apply -f micro-rpc-single-pod.yaml
sleep 180s
hey -n $N -c $C -t 1000 -d "{}" -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-single-pod-$DEST.csv
sudo kubectl delete -f micro-rpc-single-pod.yaml
sleep 180s
sudo kubectl apply -f micro-rpc-single-mmap.yaml
sleep 180s
hey -n $N -c $C -t 1000 -d "{}" -o csv -m POST -T application/json http://10.43.190.1 >> micro-rpc-single-mmap-$DEST.csv
sudo kubectl delete -f micro-rpc-single-mmap.yaml
sleep 180s
sed -i "s/$SIZE/1000000/g" ../../benchmarks/micro/micro-rpc-a/func.yaml
sed -i "s/$SIZE/1000000/g" ../../benchmarks/micro/micro-rpc-b/func.yaml
sed -i "s/$SIZE/1000000/g" ../../benchmarks/micro/micro-rpc-a-b/func.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-multi-pod.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-single-pod.yaml
sed -i "s/$SIZE/1000000/g" micro-rpc-single-mmap.yaml
