#!/bin/bash -e

set -x

service=sites

# group=public

cp ~/.ssh/authorized_keys .ssh/

# version=$(vag docker version patch ${service}-${group})

version=$(platformctl next tag 7onetella/sites)

docker build -t docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}" .

docker push docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}"

# nomad job stop -purge ${service}

# vag docker deploy docker-registry.int.curiosityworks.org/7onetella/${service}-${group}:"${version}"