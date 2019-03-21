#!/bin/bash

arch_list="amd64 arm arm64 ppc64le s390x"
manifest="${DOCKER_USER:="jeefy"}/dashboard-metrics-sidecar";
manifest_list="";


for i in ${arch_list}; do
    # If it's a tagged release, use the tag
    # Otherwise, assume it's HEAD and push to latest
    container="${manifest}-${i}:${TRAVIS_TAG:="latest"}"
    
    echo "--- docker build --no-cache --build-arg GOARCH=${i} -t ${container} .";
    docker build --no-cache --build-arg GOARCH=${i} -t ${container} .

    echo "--- docker push ${container}"
    docker push ${container}
    
    manifest_list="${manifest_list} ${container}"
done;

echo "--- docker manifest create --amend ${manifest} ${manifest_list}"
docker manifest create --amend ${manifest} ${manifest_list}

for i in ${arch_list}; do
    container="${manifest}-${i}:${TRAVIS_TAG:="latest"}"

    echo "--- docker manifest annotate ${manifest} ${container} --os linux --arch ${i}"
    docker manifest annotate ${manifest} ${container} --os linux --arch ${i}
done;

echo "--- docker manifest push ${manifest}"
docker manifest push ${manifest}