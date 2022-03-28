#!/bin/bash -e

set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

gox -osarch="linux/amd64"

version=0.1.0

docker build -t curiosinauts/console:${version} .

docker push curiosinauts/console:${version}

cat console.tpl | sed 's/__tag__/'"${version}"'/g' > console.yml

kubectl delete -f ./console.yml || true

kubectl apply -f ./console.yml

rm -f console.yml || true

rm -f console_linux_amd64 || true