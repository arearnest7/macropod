#bin/bash
host=${1:-127.0.0.1}
for i in yamls/election*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\\.$host/g" $i; done;

kubectl apply -f yamls/election-gateway.yaml 2>dev/null
kubectl apply -f yamls/election-get-results.yaml 2>dev/null
kubectl apply -f yamls/election-vote-enqueuer.yaml 2>dev/null
kubectl apply -f yamls/election-vote-processor.yaml 2>dev/null

for i in yamls/election*; do sed -i "s/knative-functions\\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
