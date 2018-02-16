#!/bin/bash
RED='\e[91m'
GREEN='\e[92m'
BLUE='\e[36m'
YELLOW='\e[93m'
RESET='\e[39m'
APP=chad
IMAGENAME=${IMAGENAME:=ppamo.cl/chad}
DOCKERIMAGE="${DOCKERIMAGE:=docker.io/golang:1.9.2-alpine}"
PROJECT=github.com/Ppamo/go.sidecar.ambassador/chad
BINARYFILE=bin/chad
BINARYDOCKER=docker/chad
VERSION=${VERSION:=0.1.0}

if [ -z "$GOPATH" ]
then
	echo "* GOPATH is not set!"
	exit -1
fi

usage(){
	printf "$YELLOW* Usage:
	build $BLUE[compile|build|run|clean|list]$RESET
	"
}

stop(){
	ID=$(docker ps --format="{{.ID}}" --filter="name=$APP")
	if [ "$ID" ]
	then
		docker stop $ID
	else
		printf "$RED- ERROR: No container found$RESET\n"
	fi
}

sigint_handler(){
	printf "$YELLOW* Stopping app $RESET\n"
	stop
}

compile(){
	printf $YELLOW"* Compiling chad$BLUE\n"
	rm -f $BINARYFILE
	OS=$(uname -o)
	if [ "$OS" == "Cygwin" ]
	then
		GOPATH=${GOPATH#/cygdrive}
	fi


	CMD="cd \"src/$PROJECT\" && \
		go get && \
		CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o \"$BINARYFILE\" ."

	docker run --rm --privileged=true -i -v "$GOPATH:/go" "$DOCKERIMAGE" /bin/sh -c "$CMD"
	if [ -x $BINARYFILE ]
	then
		printf $GREEN"Chad OK!$RESET\n"
	else
		printf $RED"Chad Not OK!$RESET\n"
	fi
}

build(){
	printf $YELLOW"* Creating chad\n"
	rm -f $BINARYDOCKER
	if [ ! -x $BINARYFILE ]
	then
		printf $RED"+ Chad not found, compile first, OK!$RESET\n"
		exit -2
	fi
	cp "$BINARYFILE" docker/
	docker images $IMAGENAME --format "{{.Tag}}" | grep "$VERSION" > /dev/null
	if [ $? -eq 0 ]
	then
		printf "$YELLOW+ Deleting image $IMAGENAME:$VERSION $RESET\n"
		docker rmi $IMAGENAME:$VERSION
	fi
	printf "$YELLOW+ Building docker image $IMAGENAME:$VERSION $RESET\n"
	docker build -t $IMAGENAME:$VERSION docker/
	if [ $? -eq 0 ]
	then
		printf $GREEN"Chad OK!$RESET\n"
	else
		printf $RED"Chad Not OK!$RESET\n"
	fi
}

run(){
	printf "$YELLOW* Getting container $IMAGENAME $RESET\n"
	VERSION=$(docker images $IMAGENAME --format "{{.Tag}}" | \
		grep -v "latest" | \
		sort | \
		tail -n 1)
	if [ "$VERSION" ]
	then
		printf "$YELLOW+ Running $IMAGENAME:$VERSION $RESET\n"
		trap sigint_handler SIGINT
		docker run -p 8081:8081 --rm --name "$APP" -i $IMAGENAME:$VERSION
	else
		printf "$RED- ERROR: no image found, please compile and build first, OK! $RESET\n"
	fi
	exit 0
}

clean(){
	printf "$YELLOW* Cleaning chad! $RESET\n"
	printf "$YELLOW+ Deleting $BINARYFILE$RESET\n"
	rm -f $BINARYFILE
	printf "$YELLOW+ Deleting $BINARYDOCKER$RESET\n"
	rm -f $BINARYDOCKER
	printf "$YELLOW+ Stopping app $RESET\n"
	stop
	TAGS=$(docker images ppamo.cl/chad --format "{{.Tag}}")
	for TAG in $TAGS
	do
		printf "$YELLOW+ Deleting $IMAGENAME:$TAG$RESET\n"
		docker rmi $IMAGENAME:$TAG
	done
	printf $GREEN"Chad OK!$RESET\n"
}

list(){
	printf "$YELLOW* Listing chad images $RESET\n"
	docker images $IMAGENAME --format "- {{.Repository}}:{{.Tag}}\t{{.CreatedAt}}\t{{.Size}}"
	printf $GREEN"Chad OK!$RESET\n"
}

case $1 in
	compile)	compile		;;
	build)		build		;;
	run)		run		;;
	clean)		clean		;;
	list)		list		;;
	*)
		printf "$RED- ERROR: comando no encontrado$RESET\n"
		usage
		exit -1
esac
