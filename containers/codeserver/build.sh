#!/bin/bash -e

username=$1 

service=vscode-${username}

set -x

codeserver_base=1.0.12

version=$2

platformctl before docker-build ${username} ${version}

docker build --no-cache --build-arg CODESERVER_BASE=${codeserver_base} -t docker-registry.int.curiosityworks.org/7onetella/${service}:${version} .

docker push docker-registry.int.curiosityworks.org/7onetella/${service}:${version}

kubectl apply -f ./vscode-${username}.yml

platformctl after docker-build ${username}

docker images | grep vscode | awk '{ print $3 }' | xargs docker rmi