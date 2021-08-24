#!/bin/bash -e

set -x

go build
echo

platformctl add user 7onetella@gmail.com
echo

platformctl list users
echo

platformctl list repos
echo

platformctl describe user 7onetella@gmail.com
echo

platformctl remove user 7onetella@gmail.com