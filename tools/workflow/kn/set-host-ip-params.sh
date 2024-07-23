#!/bin/bash
find ../../.. -name "func.yaml" -print0 | while read -d $'\0' file; do sed -i "s/knative-functions.$1/knative-functions.$2/g" $file; done;
find ./yamls -name "*.yaml" -print0 | while read -d $'\0' file; do sed -i "s/knative-functions.$1/knative-functions.$2/g" $file; done;
