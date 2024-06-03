#!/bin/bash
kn func delete -p ../../../benchmarks/$1/original/feature-generation/feature-extractor
kn func delete -p ../../../benchmarks/$1/original/feature-generation/feature-orchestrator
kn func delete -p ../../../benchmarks/$1/original/feature-generation/feature-reducer
kn func delete -p ../../../benchmarks/$1/original/feature-generation/feature-status
kn func delete -p ../../../benchmarks/$1/original/feature-generation/feature-wait

