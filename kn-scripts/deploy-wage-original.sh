#bin/bash
host=${1:-127.0.0.1}
for i in yamls/wage*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

kubectl apply -f yamls/wage-avg.yaml 2>dev/null
kubectl apply -f yamls/wage-format.yaml 2>dev/null
kubectl apply -f yamls/wage-merit.yaml 2>dev/null
kubectl apply -f yamls/wage-stats.yaml 2>dev/null
kubectl apply -f yamls/wage-sum.yaml 2>dev/null
kubectl apply -f yamls/wage-validator.yaml 2>dev/null
kubectl apply -f yamls/wage-write-merit.yaml 2>dev/null
kubectl apply -f yamls/wage-write-raw.yaml 2>dev/null

for i in yamls/wage*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
