find .. -name "func.yaml" -print0 | while read -d $'\0' file; do sed -i "s/value: $1/value: $2/g" $file; done;
find .. -name "func.yaml" -print0 | while read -d $'\0' file; do sed -i "s/value: $3/value: $4/g" $file; done;
