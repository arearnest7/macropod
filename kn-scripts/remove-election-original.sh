#bin/bash
sudo kubectl delete -f yamls/election-gateway.yaml
sudo kubectl delete -f yamls/election-get-results.yaml
sudo kubectl delete -f yamls/election-vote-enqueuer.yaml
sudo kubectl delete -f yamls/election-vote-processor.yaml
