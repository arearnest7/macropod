#!/bin/bash
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-avg
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-format
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-merit
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-stats
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-sum
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-validator
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-write-merit
kn func delete -p ../../../benchmarks/$1/original/wage-pay/wage-write-raw

