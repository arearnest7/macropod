#bin/bash
sudo kubectl delete -f yamls/wage-avg.yaml
sudo kubectl delete -f yamls/wage-format.yaml
sudo kubectl delete -f yamls/wage-merit.yaml
sudo kubectl delete -f yamls/wage-stats.yaml
sudo kubectl delete -f yamls/wage-sum.yaml
sudo kubectl delete -f yamls/wage-validator.yaml
sudo kubectl delete -f yamls/wage-write-merit.yaml
sudo kubectl delete -f yamls/wage-write-raw.yaml

