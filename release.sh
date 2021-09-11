#!/bin/bash -e

if [ "${1}" == "" ]; then
  echo "specify tag version"
  exit 1
fi

version=${1}

git tag -a v${version} -m "${version} release"

git push origin --tags

goreleaser release --rm-dist

./platformctl run jenkins-job upgrade-platformctl -p PLATFORMCTL_VERSION="${version}"