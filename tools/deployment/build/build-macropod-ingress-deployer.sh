#!/bin/bash
ID=${1:-sysdevtamu}
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:ttl-reclaim && sudo docker push $ID/macropod-ingress:ttl-reclaim
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:ttl-reclaim && sudo docker push $ID/macropod-deployer:ttl-reclaim
