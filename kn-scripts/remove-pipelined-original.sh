#bin/bash
kubectl delete -f yamls/pipelined-checksum.yaml
kubectl delete -f yamls/pipelined-encrypt.yaml
kubectl delete -f yamls/pipelined-main.yaml
kubectl delete -f yamls/pipelined-zip.yaml

