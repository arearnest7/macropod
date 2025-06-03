#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-macropod}
./build-macropod.sh $ID $TAG
./build-kn.sh $ID
