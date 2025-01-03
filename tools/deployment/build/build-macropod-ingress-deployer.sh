#!/bin/bash
ID=${1:-sysdevtamu}
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:node-reclaim && sudo docker push $ID/macropod-ingress:node-reclaim
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:node-reclaim && sudo docker push $ID/macropod-deployer:node-reclaim
