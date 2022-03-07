#!/bin/bash -e

set -x

service=sites

cp ~/.ssh/authorized_keys .ssh/

version=$(platformctl next tag curiosinauts/sites)
# version=0.1.0

docker build -t docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}" .

docker push docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}"

