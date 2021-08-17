#!/bin/bash

version=${1}

git tag -a v${version} -m "${version} release"

git push origin --tags

goreleaser release --rm-dist