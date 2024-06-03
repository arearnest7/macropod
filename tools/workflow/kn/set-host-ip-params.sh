#!/bin/bash
find ../.. -name "func.yaml" -print0 | while read -d $'\0' file; do sed -i "s/default.$1/default.$2/g" $file; done;
