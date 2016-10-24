#!/bin/bash

echo "export GOPATH"
cd ../../
echo $(pwd)
export GOPATH=$(pwd)
echo "GOPATH:"$GOPATH

echo "get packages..."
go get github.com/emicklei/go-restful
go get github.com/Sirupsen/logrus
echo "get packages finished"

echo "build..."
go build -o deployer
echo "build finished"
