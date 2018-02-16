#!/bin/bash
RED='\e[91m'
GREEN='\e[92m'
BLUE='\e[36m'
YELLOW='\e[93m'
RESET='\e[39m'
APP=chad
IMAGENAME=${IMAGENAME:=ppamo.cl/$APP}
GOIMAGE="${GOIMAGE:=docker.io/golang:1.9.2-alpine}"
PROJECT=github.com/Ppamo/go.sidecar.ambassador/chad
BINARYFILE=bin/chad
BINARYDOCKER=docker/chad
VERSION=${VERSION:=0.1.0}
KUBECTL='kubectl'
KUBEENV=$($KUBECTL config current-context)

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

load_config_file(){
	printf $YELLOW"* Cargando configuracion desde $1\n"$RESET
	if [ ! -f "$1" ]
	then
		printf $RED"- ERROR:"$RESET" Archivo de configuracion \"$1\" no existe\n"
		exit -1
	fi
	for i in $(cat $1)
	do
		VARNAME=$(echo $i | awk -F= '{ print $1 }')
		VARVALUE=$(echo $i | awk -F= '{ print $2 }')
		if [ -z "${!VARNAME}" ]
		then
			export ${VARNAME}=$VARVALUE
			printf "$GREEN+ $VARNAME"$RESET"=$VARVALUE\n"
		fi
	done
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

	docker run --rm --privileged=true -i -v "$GOPATH:/go" "$GOIMAGE" /bin/sh -c "$CMD"
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

push(){
	DST=$2
	if [ -z "$VERSION" ]
	then
		VERSION=$(docker images $IMAGENAME --format "{{.Tag}}" | \
			grep -v "latest" | \
			sort | \
			tail -n 1)
	fi
	printf "$YELLOW* Pushing chad image $IMAGENAME:$VERSION $RESET\n"
	if [ -z "$VERSION" ]
	then
		printf $RED"- ERROR: No version found!\n"$RESET
		exit -1
	fi

	docker tag $IMAGENAME:$VERSION $DST/$APP:$VERSION && \
	docker push $DST/$APP:$VERSION
	if [ $? -eq 0 ]
	then
		printf $GREEN"Chad OK!$RESET\n"
	else
		printf $RED"Chad Not OK!$RESET\n"
	fi
	((INDEX++))
}

deploy(){
	REGISTRYHOST=$2
	DEPLOYMENTPROPERTIES=deploy/app.properties
	DEPLOYMENTTEMPLATE=deploy/template.deployment
	printf "$YELLOW* Deploying chad $RESET\n"
	load_config_file $DEPLOYMENTPROPERTIES
	if [ -z "$VERSION" ]
	then
		VERSION=$(docker images $IMAGENAME --format "{{.Tag}}" | \
			grep -v "latest" | \
			sort | \
			tail -n 1)
	fi
	$KUBECTL get namespaces > /dev/null 2>&1
	if [ $? -ne 0 ]
	then
		printf "$RED- ERROR: no se pudo acceder a kubernetes en$GREEN $KUBEENV\n"$RESET
		exit -1
	fi

	printf $YELLOW"* Creando archivo deployment\n"$RESET
	IFS=$'\n'
	while read line
	do
		eval echo \"$line\"
	done < $DEPLOYMENTTEMPLATE > $DEPLOYMENTTEMPLATE.yaml
	unset IFS

	printf $YELLOW"* Desplegando proyecto$GREEN $PROJECTNAME$YELLOW a contexto$GREEN $KUBEENV\n"$RESET
	$KUBECTL get deploy "$PROJECTNAME" --namespace "$NAMESPACE" > /dev/null 2>&1
	if [ $? -eq 0 ]
	then
		$KUBECTL delete -f $DEPLOYMENTTEMPLATE.yaml
		if [ $? -ne 0 ]
		then
			printf "$RED- ERROR: no se pudo eliminar el proyecto\n"$RESET
			exit -1
		fi
	fi
	$KUBECTL create --validate=false -f $DEPLOYMENTTEMPLATE.yaml
	if [ $? -ne 0 ]
	then
		printf "$RED- ERROR: no se pudo crear el proyecto\n"$RESET
		exit -1
	fi
	printf $GREEN"* Hecho!\n"$RESET
	((INDEX++))
}

INDEX=0
while [ $INDEX -lt $# ]
do
	((INDEX++))
	case ${@:INDEX:1} in
		compile)	compile			;;
		build)		build			;;
		run)		run			;;
		clean)		clean			;;
		list)		list			;;
		push)		push ${@:INDEX}		;;
		deploy)		deploy ${@:INDEX}	;;
		*)
			printf $RED"- ERROR: comando \"${@:INDEX:1}\" no reconocido\n"$RESET
			usage
			exit -1
	esac
done
