#!/bin/bash

if [ -z "$GOPATH" ]
then
	echo "* GOPATH is not set!"
	exit -1
fi

# functions
sigint_handler(){
	printf "${ORANGE}* Stopping app ${NC}\n"
	ID=$(docker ps --format="{{.ID}}" --filter="name=$APP")
	if [ "$ID" ]
	then
		docker stop $ID
	else
		echo "no container found"
	fi
}

# constants
APP="sidecar.ambassador"
VERSION=0.1.0
SRC="github.com/Ppamo/go.sidecar.ambassador"
DST="docker/$APP"
IMAGENAME="ppamo.cl/$APP"
RED='\033[0;31m'
BLUE='\033[0;34m'
ORANGE='\033[0;33m'
NC='\033[0m'

# vars
BINARYFILE="$GOPATH/src/$SRC/$DST"

# clean up
rm -f $BINARYFILE

# build
printf "${ORANGE}* Building app ${NC}\n"
bash build.sh "$SRC" "$DST"

if [ -x $BINARYFILE ]
then
	IMAGEVERSION=$(docker images $IMAGENAME --format "{{.Tag}}")
	if [ "$IMAGEVERSION" ]
	then
		printf "${ORANGE}* Deleting image $IMAGENAME:$IMAGEVERSION ${NC}\n"
		docker rmi $IMAGENAME:$IMAGEVERSION
	fi
	printf "${ORANGE}* Building docker image $IMAGENAME:$VERSION ${NC}\n"
	docker build -t $IMAGENAME:$VERSION docker/

	IMAGESNUMBER=$(docker images $IMAGENAME:$VERSION --format "{{.ID}}" | wc -l)
	if [ $IMAGESNUMBER -gt 0 ]
	then
		printf "${ORANGE}* Running container ${NC}\n"
		trap sigint_handler SIGINT
		docker run --rm --name "$APP" -i $IMAGENAME:$VERSION
	fi
	exit 0
fi

printf  "${RED}* Build failed!!${NC}\n"
exit 1
