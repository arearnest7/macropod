#!/bin/bash
ID=${1:-arearnest7}
docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:latest && docker push $ID/macropod-ingress:latest
docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:latest && docker push $ID/macropod-deployer:latest
docker build ../../../macropod-base/macropod-prepuller -t $ID/macropod-prepuller:latest && docker push $ID/macropod-prepuller:latest
