#!/bin/bash

echo 'Deploying InvisibleURL.net'

if [ -z "$1" ] || [ -z "$2" ]
then
    echo 'No argument supplied'
    exit
fi

echo "Deploying using the ssh key: " $1 " to the new release folder, " $2

# run all the tests. If one of them fails, do not deploy
if ! go test ./...;
then
    echo "There were a test failure. Aborting deploy"
    exit
else
    echo "All tests passed successfully."
fi

# build a release, copy the assets to the release folder
echo "Starting build"
if ! env GOOS=linux GOARCH=arm64 go build -v -o ./release;
then
    echo "Build failed. Aborting deploy";
    exit
else
    echo "Build succeeded"
fi

echo "Copying static assets"
cp -r ./public/ ./release/
cp -r ./views/ ./release/
cp start.sh ./release/

# create new release folder and copy the files to the server
echo "Copying release to the server"
if ! ssh -i $1 invisibleurl@116.203.110.44 "mkdir -p /opt/invisibleurl/releases/$2"
then
    echo "Couldn't create the release directory on the server"
    exit
fi

if ! scp -i $1 -r ./release/* invisibleurl@116.203.110.44:/opt/invisibleurl/releases/$2
then
    echo "Couldn't copy the release files to the server"
    exit
fi
echo "Copying release to the server finished"
echo "Please finish the release on the server"

