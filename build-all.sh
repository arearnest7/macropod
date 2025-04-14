#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
./build-macropod.sh $ID $TAG
./build-kn.sh $ID
