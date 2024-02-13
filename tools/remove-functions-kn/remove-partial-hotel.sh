#!/bin/bash
PATH=${1:-kn}
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-frontend-spgr
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-recommend-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-reserve-partial
kn func delete -p ../../benchmarks/$PATH/full-reduced/serverless-election/hotel-user-partial
