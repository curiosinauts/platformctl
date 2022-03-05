#!/bin/bash -e

set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

repo=languages

base_version=${1}

version=$(platformctl next tag 7onetella/languages)

docker build --build-arg BASE_VERSION=${base_version} -t docker-registry.int.curiosityworks.org/7onetella/${repo}:"${version}" .

docker push docker-registry.int.curiosityworks.org/7onetella/${repo}:"${version}"
