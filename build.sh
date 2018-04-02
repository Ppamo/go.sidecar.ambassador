#!/bin/bash
APP="sidecar.ambassador"
PROJECT=github.com/Ppamo/go.sidecar.ambassador
source "$(git rev-parse --show-toplevel)/lib/build_lib.sh"

usage(){
	printf "$YELLOW* Usage:
	build $BLUE[compile|build|run|clean|list]$RESET
"
}

if [ $# -eq 0 ]
then
	usage
	exit 0
fi

INDEX=0
while [ $INDEX -lt $# ]
do
	((INDEX++))
	case ${@:INDEX:1} in
		compile)	compile ${@:INDEX}	;;
		goget)		goget ${@:INDEX}	;;
		build)
			DOCKERCOPYFILES="config/config.json" \
			build
			;;
		run)
			DUMMYID=$(docker ps --filter name=dummy --format "{{.ID}}")
			if [ $DUMMYID ]
			then
				DUMMYIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $DUMMYID)
			fi
			DOCKERENV="-e SERVERPORT=8080 -e DESTINATION=http://$DUMMYIP:8081" \
			run
			;;
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
