#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
sudo docker buildx build base/macropod-ingress -t $ID/macropod-ingress:$TAG && sudo docker push $ID/macropod-ingress:$TAG
sudo docker buildx build base/macropod-deployer -t $ID/macropod-deployer:$TAG && sudo docker push $ID/macropod-deployer:$TAG

