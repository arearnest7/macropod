#!/bin/bash
sudo kubectl delete -f ../yamls/sn-compose-and-build.yaml
sudo kubectl delete -f ../yamls/sn-compose-post.yaml
sudo kubectl delete -f ../yamls/sn-db-op.yaml
sudo kubectl delete -f ../yamls/sn-entry.yaml
sudo kubectl delete -f ../yamls/sn-post-storage.yaml
sudo kubectl delete -f ../yamls/sn-upload-creator.yaml
sudo kubectl delete -f ../yamls/sn-upload-home-timeline.yaml
sudo kubectl delete -f ../yamls/sn-upload-media.yaml
sudo kubectl delete -f ../yamls/sn-upload-text.yaml
sudo kubectl delete -f ../yamls/sn-upload-unique-id.yaml
sudo kubectl delete -f ../yamls/sn-upload-user-mentions.yaml
sudo kubectl delete -f ../yamls/sn-upload-user-timeline.yaml
