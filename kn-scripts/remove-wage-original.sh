#bin/bash
kubectl delete -f yamls/wage-avg.yaml
kubectl delete -f yamls/wage-format.yaml
kubectl delete -f yamls/wage-merit.yaml
kubectl delete -f yamls/wage-stats.yaml
kubectl delete -f yamls/wage-sum.yaml
kubectl delete -f yamls/wage-validator.yaml
kubectl delete -f yamls/wage-write-merit.yaml
kubectl delete -f yamls/wage-write-raw.yaml

