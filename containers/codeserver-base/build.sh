#!/bin/bash -e

set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

service=codeserver-base

tools_version=${1}

version=$(platformctl next tag 7onetella/codeserver-base)

# removed --no-cache to improve build time
docker build --build-arg TOOLS_VERSION=${tools_version} -t docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}" .

docker push docker-registry.int.curiosityworks.org/7onetella/${service}:"${version}"


# chnagelog
# 1.0.1 - reverting code server back to 3.10.0 to fix python plugin click through issue
# 1.0.2 - adding sorenisanerd/gotty web terminal app