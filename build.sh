#!/bin/bash
DOCKERIMAGE="${DOCKERIMAGE:-docker.io/golang:1.9.2-alpine}"
PROJECT=$1
BINARYFILE=$2

OS=$(uname -o)
if [ "$OS" == "Cygwin" ]
then
	GOPATH=${GOPATH#/cygdrive}
fi

CMD="cd \"src/$PROJECT\" && \
	go get && \
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o \"$BINARYFILE\" ."

docker run --rm --privileged=true -i -v "$GOPATH:/go" "$DOCKERIMAGE" /bin/sh -c "$CMD"
