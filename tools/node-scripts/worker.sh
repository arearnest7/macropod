#!/bin/bash
MASTER_IP=$1
TOKEN=$2
curl -sfL https://get.k3s.io | K3S_URL=https://$MASTER_IP:6443 K3S_TOKEN=$TOKEN sh -
