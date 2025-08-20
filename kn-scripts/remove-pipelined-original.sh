#bin/bash
sudo kubectl delete -f yamls/pipelined-checksum.yaml
sudo kubectl delete -f yamls/pipelined-encrypt.yaml
sudo kubectl delete -f yamls/pipelined-main.yaml
sudo kubectl delete -f yamls/pipelined-zip.yaml

