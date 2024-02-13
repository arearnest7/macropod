#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-frontend
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-geo
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-profile
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-rate
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-recommend
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-reserve
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-search
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-user
