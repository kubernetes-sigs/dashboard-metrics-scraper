#!/bin/bash
for i in amd64 arm armv7; do
    container="metrics-sidecar-${i}"
    echo "Now pushing ${container}"
    docker push --build-arg GOARCH=${i} -t ${container} .
done;