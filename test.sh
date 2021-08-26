#!/bin/bash -e

go build

echo platformctl add user 7onetella@gmail.com
platformctl add user 7onetella@gmail.com

echo platformctl list users
platformctl list users

echo platformctl list repos
platformctl list repos

echo platformctl describe user 7onetella@gmail.com
platformctl describe user 7onetella@gmail.com

echo platformctl update codeserver 7onetella@gmail.com
platformctl update codeserver 7onetella@gmail.com

echo platformctl remove user 7onetella@gmail.com
platformctl remove user 7onetella@gmail.com