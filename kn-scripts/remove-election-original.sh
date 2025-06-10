#bin/bash
kubectl delete -f yamls/election-gateway.yaml
kubectl delete -f yamls/election-get-results.yaml
kubectl delete -f yamls/election-vote-enqueuer.yaml
kubectl delete -f yamls/election-vote-processor.yaml
