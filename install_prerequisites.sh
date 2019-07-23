#!/bin/bash

export GOPATH=`pwd`

go get github.com/docker/docker/client
go get github.com/aws/aws-sdk-go/...
go get -u github.com/Azure/azure-sdk-for-go/...
go get github.com/grokify/html-strip-tags-go
