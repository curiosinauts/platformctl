#!/bin/bash -e

set -x

go build

platformctl add user 7onetella@gmail.com

platformctl list users

platformctl list repos

platformctl describe user 7onetella@gmail.com

platformctl remove user 7onetella@gmail.com