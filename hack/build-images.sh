#!/bin/bash
for i in amd64 arm armv7; do
    container="metrics-sidecar-${i}"
    echo "Now building ${container}"
    docker build --build-arg GOARCH=${i} -t ${container} .
done;