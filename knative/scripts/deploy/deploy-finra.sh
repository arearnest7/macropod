#!/bin/bash
host=${1:-127.0.0.1}
for i in ../yamls/finra*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

sudo kubectl apply -f ../yamls/finra-entry.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-lastpx.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-margin-balance.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-marketdata.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-side.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-trddate.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/finra-volume.yaml 2>/dev/null

for i in ../yamls/finra*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
