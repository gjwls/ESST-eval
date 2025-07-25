#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <Latency(ms)> <number of processes>"
    exit 1
fi

latency=$1
process_count=$2

for i in $(seq 1 $process_count); do
    ./eval_ls $latency 1 &
done

wait
