#!/bin/bash
kn func delete -p ../../benchmarks/$1/partial-reduced/hotel-app/hotel-frontend-spgr
kn func delete -p ../../benchmarks/$1/partial-reduced/hotel-app/hotel-recommend-partial
kn func delete -p ../../benchmarks/$1/partial-reduced/hotel-app/hotel-reserve-partial
kn func delete -p ../../benchmarks/$1/partial-reduced/hotel-app/hotel-user-partial
