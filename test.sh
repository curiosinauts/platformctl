#!/bin/bash -e

go build

user_email=${1}

echo platformctl add user ${user_email} -p -e 
platformctl add user ${user_email} -p -e
echo -------------------------------------------------------------------

echo platformctl add ide-repo ${user_email} vscode -r git@github.com:curiosinauts/platformctl.git 
platformctl add ide-repo ${user_email} vscode -r git@github.com:curiosinauts/platformctl.git 
echo -------------------------------------------------------------------

echo platformctl list ide-repo ${user_email} vscode 
platformctl list ide-repo ${user_email} vscode 
echo -------------------------------------------------------------------

echo platformctl remove ide-repo ${user_email} vscode -r git@github.com:curiosinauts/platformctl.git 
platformctl remove ide-repo ${user_email} vscode -r git@github.com:curiosinauts/platformctl.git 
echo -------------------------------------------------------------------

echo platformctl list ide-repo ${user_email} vscode 
platformctl list ide-repo ${user_email} vscode 
echo -------------------------------------------------------------------

echo platformctl describe user ${user_email}
platformctl describe user ${user_email}
echo -------------------------------------------------------------------

echo platformctl list users
platformctl list users
echo -------------------------------------------------------------------

echo platformctl add runtime-install ${user_email} vscode poetry -n 
platformctl add runtime-install ${user_email} vscode poetry -n
echo -------------------------------------------------------------------

echo platformctl list users
platformctl list users
echo -------------------------------------------------------------------

echo platformctl update code-server ${user_email}  
platformctl update code-server ${user_email} 
echo -------------------------------------------------------------------

echo platformctl list repos
platformctl list repos
echo -------------------------------------------------------------------

echo platformctl remove user ${user_email}
platformctl remove user ${user_email}
echo -------------------------------------------------------------------

echo platformctl list repos
platformctl list repos
echo ------------------------------------------------------------------