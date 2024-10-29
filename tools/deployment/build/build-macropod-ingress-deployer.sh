#!/bin/bash
ID=${1:-sysdevtamu}
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:latest && sudo docker push $ID/macropod-ingress:latest
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:latest && sudo docker push $ID/macropod-deployer:latest
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-prepuller -t $ID/macropod-prepuller:latest && sudo docker push $ID/macropod-prepuller:latest
