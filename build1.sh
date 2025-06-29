#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
sudo docker buildx build base/macropod-node -t $ID/macropod-node:$TAG && sudo docker push $ID/macropod-node:$TAG

sudo docker buildx build benchmarks/macropod/unified/serverless-election/election-unified -t $ID/election-unified:$TAG && sudo docker push $ID/election-unified:$TAG

sudo docker buildx build benchmarks/macropod/original/serverless-election/election-gateway -t $ID/election-gateway:$TAG && sudo docker push $ID/election-gateway:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-get-results -t $ID/election-get-results:$TAG && sudo docker push $ID/election-get-results:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-vote-enqueuer -t $ID/election-vote-enqueuer:$TAG && sudo docker push $ID/election-vote-enqueuer:$TAG
sudo docker buildx build benchmarks/macropod/original/serverless-election/election-vote-processor -t $ID/election-vote-processor:$TAG && sudo docker push $ID/election-vote-processor:$TAG
