#!/bin/sh
# entrypoint.sh

ulimit -l unlimited

exec python index.py