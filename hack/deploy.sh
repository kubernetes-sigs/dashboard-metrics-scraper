#!/bin/bash

arch_list="amd64 arm arm64 ppc64le s390x"
manifest="${DOCKER_USER:="jeefy"}/dashboard-metrics-sidecar";
manifest_list="";


for i in ${arch_list}; do
    container="${manifest}-${i}"
    
    echo "Now building ${container}"
    docker build --build-arg GOARCH=${i} -t ${container} .

    echo "Now pushing ${container}"
    docker push ${container}
    
    manifest_list="${manifest_list} ${manifest}-${i}"
done;

docker manifest create --amend ${manifest} ${manifest_list}

for i in ${arch_list}; do
    docker manifest annotate ${manifest} "${manifest}-${i}" --os linux --arch ${i}
done;

docker manifest push ${manifest}