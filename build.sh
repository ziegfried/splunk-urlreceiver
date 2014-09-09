#!/bin/bash

APP_NAME=urlreceiver

cd $(dirname $0)

mkdir -p src pkg bin dist tmp

export GOPATH=$PWD
export GOBIN=$PWD/bin
export PATH="$GOBIN:$PATH"

echo "Cleaning staging environment..." >&2
rm -rf tmp/*

echo "Bootstrapping staging environment..." >&2
mkdir tmp/$APP_NAME
cp -r app/* tmp/$APP_NAME/

echo "Building binaries..." >&2

function build_binary() {
    echo "Building $3..." >&2
    GOOS=$1 GOARCH=$2 go build -o tmp/$APP_NAME/$3/bin/$4 src/urlreceiver.go
}

build_binary darwin amd64 darwin_x86_64 urlreceiver
build_binary darwin 386 darwin_x86 urlreceiver
build_binary linux amd64 linux_x86_64 urlreceiver
build_binary linux 386 linux_x86 urlreceiver
build_binary windows amd64 windows_x86_64 urlreceiver.exe
build_binary windows 386 windows_x86 urlreceiver.exe

echo "Building pacakge..." >&2
VERSION=$(cat app/default/app.conf | grep -e 'version\s*=' | sed -e 's/version[ ]*=[ ]*//')
FILENAME=urlreceiver-$VERSION-$(git rev-parse --short HEAD).spl

cd tmp
tar cvfz ../dist/$FILENAME urlreceiver >&2

echo >&2
echo $FILENAME
