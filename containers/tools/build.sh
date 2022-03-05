#!/bin/bash -e

set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

languages_version=${1}

version=$(platformctl next tag 7onetella/tools)

service=tools

docker build --build-arg LANGUAGES_VERSION=${languages_version} -t docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}" .

docker push docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}"
