#!/bin/bash
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-frontend
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-geo
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-profile
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-rate
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-recommend
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-reserve
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-search
kn func delete -p ../../benchmarks/$1/original/hotel-app/hotel-user
