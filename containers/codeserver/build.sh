#!/bin/bash -e

username=$1 

service=vscode-${username}

set -x

base=0.1.1

version=$2

platformctl before docker-build ${username} ${version}

docker build --no-cache --build-arg BASE=${base} -t docker-registry.curiosityworks.org/curiosinauts/${service}:${version} .

docker push docker-registry.curiosityworks.org/curiosinauts/${service}:${version}

kubectl apply -f ./vscode-${username}.yml

platformctl after docker-build ${username}

docker images | grep vscode | awk '{ print $3 }' | xargs docker rmi