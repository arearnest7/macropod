#!/bin/bash
host=${1:-127.0.0.1}
for i in ../yamls/sn*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

sudo kubectl apply -f ../yamls/sn-compose-and-build.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-compose-post.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-db-op.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-entry.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-post-storage.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-creator.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-home-timeline.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-media.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-text.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-unique-id.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-user-mentions.yaml 2>/dev/null
sudo kubectl apply -f ../yamls/sn-upload-user-timeline.yaml 2>/dev/null

for i in ../yamls/sn*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
