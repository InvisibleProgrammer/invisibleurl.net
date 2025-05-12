#!/bin/bash

if [ -z "$1" ]
then
    echo 'No argument supplied'
    exit
fi

echo "Deploying the release $1"

rm -rf /opt/invisibleurl/releases/current/*
cp -r /opt/invisibleurl/releases/$1/* /opt/invisibleurl/releases/current/
chmod +x /opt/invisibleurl/releases/current/start.sh

sudo systemctl restart invisibleurl

echo "$1 released"
