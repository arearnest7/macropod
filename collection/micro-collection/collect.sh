#!/bin/bash
DEST=${2:-results}
N=${3:-10000}
C=${4:-1}
./micro-invoke.sh 1MB 1000000 $1 $N $C
./micro-invoke.sh 10MB 10000000 $1 $N $C
./micro-invoke.sh 20MB 20000000 $1 $N $C
./micro-invoke.sh 50MB 50000000 $1 $N $C
./micro-invoke.sh 100MB 100000000 $1 $N $C
mkdir $DEST
mv *.csv $DEST
