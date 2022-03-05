#!/bin/bash

cd www.7onetella.net

hugo

cd ../www.curiosityworks.org

hugo

git add .

git commit -m "publishing to live"

git push origin