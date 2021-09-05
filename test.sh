#!/bin/bash -e

go build

echo platformctl add user 7onetella@gmail.com -p -e 
platformctl add user 7onetella@gmail.com -p -e
echo -------------------------------------------------------------------

echo platformctl add ide-repo 7onetella@gmail.com vscode -r git@github.com:curiosinauts/platformctl.git 
platformctl add ide-repo 7onetella@gmail.com vscode -r git@github.com:curiosinauts/platformctl.git 
echo -------------------------------------------------------------------

echo platformctl list ide-repo 7onetella@gmail.com vscode 
platformctl list ide-repo 7onetella@gmail.com vscode 
echo -------------------------------------------------------------------

echo platformctl remove ide-repo 7onetella@gmail.com vscode -r git@github.com:curiosinauts/platformctl.git 
platformctl remove ide-repo 7onetella@gmail.com vscode -r git@github.com:curiosinauts/platformctl.git 
echo -------------------------------------------------------------------

echo platformctl list ide-repo 7onetella@gmail.com vscode 
platformctl list ide-repo 7onetella@gmail.com vscode 
echo -------------------------------------------------------------------

echo platformctl describe user 7onetella@gmail.com
platformctl describe user 7onetella@gmail.com
echo -------------------------------------------------------------------

echo platformctl list users
platformctl list users
echo -------------------------------------------------------------------

echo platformctl add runtime-install 7onetella@gmail.com vscode poetry -n 
platformctl add runtime-install 7onetella@gmail.com vscode poetry -n
echo -------------------------------------------------------------------

echo platformctl list users
platformctl list users
echo -------------------------------------------------------------------

echo platformctl update code-server 7onetella@gmail.com  
platformctl update code-server 7onetella@gmail.com 
echo -------------------------------------------------------------------

echo platformctl list repos
platformctl list repos
echo -------------------------------------------------------------------

echo platformctl remove user 7onetella@gmail.com
platformctl remove user 7onetella@gmail.com
echo -------------------------------------------------------------------

echo platformctl list repos
platformctl list repos
echo ------------------------------------------------------------------