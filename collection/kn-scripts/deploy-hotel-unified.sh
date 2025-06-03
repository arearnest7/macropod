#bin/bash
host=${1:-127.0.0.1}
for i in yamls/hotel*; do sed -i "s/knative-functions\.127\.0\.0\.1/knative-functions\.$host/g" $i; done;

kubectl apply -f yamls/hotel-unified.yaml 2>dev/null

for i in yamls/hotel*; do sed -i "s/knative-functions\.$host/knative-functions\.127\.0\.0\.1/g" $i; done;
