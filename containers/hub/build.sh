#!/bin/bash -e

set -x

service=docker-hub

group=backend

version=$(vag docker version patch ${service}-${group})

docker build -t docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}" .

docker push docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}"

vag docker deploy docker-registry.curiosityworks.org/curiosinauts/${service}-${group}:"${version}"
