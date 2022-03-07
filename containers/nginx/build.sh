#!/bin/bash -e

set -x

service=sites

# group=public

cp ~/.ssh/authorized_keys .ssh/

# version=$(vag docker version patch ${service}-${group})

version=$(platformctl next tag 7onetella/sites)

docker build -t docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}" .

docker push docker-registry.curiosityworks.org/curiosinauts/${service}:"${version}"

# nomad job stop -purge ${service}
