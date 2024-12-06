#!/bin/bash
ID=${1:-sysdevtamu}
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:http && sudo docker push $ID/macropod-ingress:http
sudo DOCKER_BUILDKIT=1 docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:http && sudo docker push $ID/macropod-deployer:http
