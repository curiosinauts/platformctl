#!/bin/bash -e

set -x

service=docker-hub

group=backend

version=$(vag docker version patch ${service}-${group})

docker build -t docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}" .

docker push docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}"

vag docker deploy docker-registry.int.curiosityworks.org/7onetella/${service}-${group}:"${version}"
