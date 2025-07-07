#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
sudo docker buildx build base/macropod-eval -t $ID/macropod-eval:$TAG && sudo docker push $ID/macropod-eval:$TAG

