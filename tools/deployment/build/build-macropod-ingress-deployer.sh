#!/bin/bash
ID=${1:-arearnest7}
sudo docker build ../../../macropod-base/macropod-ingress -t $ID/macropod-ingress:latest && sudo docker push $ID/macropod-ingress:latest
sudo docker build ../../../macropod-base/macropod-deployer -t $ID/macropod-deployer:latest && sudo docker push $ID/macropod-deployer:latest
sudo docker build ../../../macropod-base/macropod-prepuller -t $ID/macropod-prepuller:latest && sudo docker push $ID/macropod-prepuller:latest
