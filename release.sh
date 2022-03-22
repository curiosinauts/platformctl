#!/bin/bash -e

# make sure to do this before releasing 
# export GITHUB_TOKEN=ghp_xxxxxxx <= github/7onetella Settings/Developer settings/personal access token
# gorelease

if [ "${1}" == "" ]; then
  echo "specify tag version"
  exit 1
fi

version=${1}

git tag -a v${version} -m "${version} release"

git push origin --tags

goreleaser release --rm-dist

# ./platformctl run jenkins-job upgrade-platformctl -p PLATFORMCTL_VERSION="${version}"