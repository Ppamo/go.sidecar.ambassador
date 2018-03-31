#!/bin/bash
if [ -z "$GOPATH" ]; then
	echo "* GOPATH is not set!"
	exit -1
fi
if [ -z "$APP" ]; then
	echo "* APP is not set!"
	exit -1
fi

# constants & vars
KUBECTL='kubectl'
RED='\e[91m'
GREEN='\e[92m'
BLUE='\e[36m'
YELLOW='\e[93m'
RESET='\e[39m'
PROJECT=${PROJECT:=github.com/Ppamo/go.sidecar.ambassador/$APP}
IMAGENAME=${IMAGENAME:=ppamo.cl/$APP}
GOIMAGE="${GOIMAGE:=docker.io/golang:1.9.2-alpine}"
BINARYFILE=bin/$APP
BINARYDOCKER=docker/$APP
VERSION=${VERSION:=0.1.0}
PORT=${PORT:=8080}
KUBEENV=$($KUBECTL config current-context)

# logic
stop_container(){
	ID=$(docker ps --format="{{.ID}}" --filter="name=$APP")
	if [ "$ID" ]
	then
		docker stop $ID
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
	stop_container
}

compile(){
	printf $YELLOW"* Compiling $APP$BLUE\n"
	rm -f $BINARYFILE
	OS=$(uname -o)
	if [ "$OS" == "Cygwin" ]
	then
		GOPATH=${GOPATH#/cygdrive}
	fi


	CMD="cd \"src/$PROJECT\" && \
		GOROOT=/usr/local/go \
		CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o \"$BINARYFILE\" ."

	docker run --rm --privileged=true -i \
		-v "$GOPATH:/go" "$GOIMAGE" /bin/sh -c "$CMD"
	if [ -x $BINARYFILE ]
	then
		printf $GREEN"$APP OK!$RESET\n"
	else
		printf $RED"$APP Not OK!$RESET\n"
		exit -1
	fi
}

goget(){
	printf $YELLOW"* Getting $2\n$RESET"
	which git > /dev/null 2>&1
	if [ $? -ne 0 ]
	then
		printf $RED"Git not found! Not OK!$RESET\n"
		exit -1
	fi
	((INDEX++))
	LIBPATH=$GOPATH/src
	mkdir -p $LIBPATH/$2
	printf "$BLUE"
	git clone https://$2 $LIBPATH/$2
	printf "$RESET"
	if [ $? -eq 0 ]
	then
		printf $GREEN"$APP OK!$RESET\n"
	else
		printf $RED"$APP Not OK!$RESET\n"
		exit -1
	fi
}

build(){
	printf $YELLOW"* Creating $APP\n"
	rm -f $BINARYDOCKER
	if [ ! -x $BINARYFILE ]
	then
		printf $RED"+ $APP not found, compile first, OK!$RESET\n"
		exit -2
	fi
	cp "$BINARYFILE" $BINARYDOCKER
	if [ "$DOCKERCOPYFILES" ]
	then
		cp --force $DOCKERCOPYFILES docker/
	fi
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
		printf $GREEN"$APP OK!$RESET\n"
	else
		printf $RED"$APP Not OK!$RESET\n"
		exit -1
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
		docker run -p $PORT:$PORT $DOCKERENV --rm --name "$APP" -i $IMAGENAME:$VERSION
	else
		printf "$RED- ERROR: no image found, please compile and build first, OK! $RESET\n"
	fi
	exit 0
}

launch(){
	printf "$YELLOW* Getting container $IMAGENAME $RESET\n"
	VERSION=$(docker images $IMAGENAME --format "{{.Tag}}" | \
		grep -v "latest" | \
		sort | \
		tail -n 1)
	if [ "$VERSION" ]
	then
		printf "$YELLOW+ Launching $IMAGENAME:$VERSION $RESET\n"
		docker run -d -p $PORT:$PORT --name "$APP" -i $IMAGENAME:$VERSION > /dev/null
		if [ $? -ne 0 ]
		then
			printf "$RED- ERROR: fail to launch $APP, OK! $RESET\n"
			exit -1
		fi
	else
		printf "$RED- ERROR: no image found, please compile and build first, OK! $RESET\n"
	fi
}

clean(){
	printf "$YELLOW* Cleaning $APP! $RESET\n"
	printf "$YELLOW+ Deleting $BINARYFILE$RESET\n"
	rm -f $BINARYFILE
	printf "$YELLOW+ Deleting $BINARYDOCKER$RESET\n"
	rm -f $BINARYDOCKER
	printf "$YELLOW+ Stopping app $RESET\n"
	stop_container
	TAGS=$(docker images ppamo.cl/$APP --format "{{.Tag}}")
	for TAG in $TAGS
	do
		printf "$YELLOW+ Deleting $IMAGENAME:$TAG$RESET\n"
		docker rmi $IMAGENAME:$TAG
	done
	printf $GREEN"$APP OK!$RESET\n"
}

list(){
	printf "$YELLOW* Listing $APP images $RESET\n"
	docker images $IMAGENAME --format "- {{printf \"%-35s\" (printf \"%.34s\" (printf \"%s:%s\" .Repository .Tag))}}{{printf \"%.20s\" .CreatedAt}}\t{{.Size}}"
	printf $GREEN"$APP OK!$RESET\n"
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
	printf "$YELLOW* Pushing $APP image $IMAGENAME:$VERSION $RESET\n"
	if [ -z "$VERSION" ]
	then
		printf $RED"- ERROR: No version found!\n"$RESET
		exit -1
	fi

	docker tag $IMAGENAME:$VERSION $DST/$APP:$VERSION && \
	docker push $DST/$APP:$VERSION
	if [ $? -eq 0 ]
	then
		printf $GREEN"$APP OK!$RESET\n"
	else
		printf $RED"$APP Not OK!$RESET\n"
	fi
	((INDEX++))
}

deploy(){
	REGISTRYHOST=$2
	DEPLOYMENTPROPERTIES=deploy/app.properties
	DEPLOYMENTTEMPLATE=deploy/template.deployment
	printf "$YELLOW* Deploying $APP $RESET\n"
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
		$KUBECTL delete --ignore-not-found=true -f $DEPLOYMENTTEMPLATE.yaml
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
