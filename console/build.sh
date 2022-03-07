#!/bin/bash -e

set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

gox -osarch="linux/amd64"

base_version=1.0.9

version=$(git log -1 --pretty=%h)

service=console

docker build --build-arg BASE_VERSION=${base_version} -t docker-registry.curiosityworks.org/curiosinauts/${service}:${version} .

docker push docker-registry.curiosityworks.org/curiosinauts/${service}:${version}

cat console.tpl | sed 's/__tag__/'"${version}"'/g' > console.yml

kubectl delete -f ./console.yml || true

kubectl apply -f ./console.yml

rm -f console.yml || true

rm -f console_linux_amd64 || true