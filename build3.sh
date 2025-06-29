#!/bin/bash
ID=${1:-sysdevtamu}
TAG=${2:-latest}
sudo docker buildx build benchmarks/macropod/unified/hotel-app/hotel-unified -t $ID/hotel-unified:$TAG && sudo docker push $ID/hotel-unified:$TAG

sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-frontend -t $ID/hotel-frontend:$TAG && sudo docker push $ID/hotel-frontend:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-geo -t $ID/hotel-geo:$TAG && sudo docker push $ID/hotel-geo:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-profile -t $ID/hotel-profile:$TAG && sudo docker push $ID/hotel-profile:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-rate -t $ID/hotel-rate:$TAG && sudo docker push $ID/hotel-rate:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-recommend -t $ID/hotel-recommend:$TAG && sudo docker push $ID/hotel-recommend:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-reserve -t $ID/hotel-reserve:$TAG && sudo docker push $ID/hotel-reserve:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-search -t $ID/hotel-search:$TAG && sudo docker push $ID/hotel-search:$TAG
sudo docker buildx build benchmarks/macropod/original/hotel-app/hotel-user -t $ID/hotel-user:$TAG && sudo docker push $ID/hotel-user:$TAG
