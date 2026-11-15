#!/bin/bash
host=${1:-127.0.0.1}
for i in ../yamls/excamera*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

sudo kubectl apply -f ../yamls/excamera-entry.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-group.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-rebase.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-reencode.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-vpxenc.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-vx-con.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/excamera-xcdec.yaml 2>/dev/null

for i in ../yamls/excamera*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
