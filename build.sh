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
		compile)	compile			;;
		build)
			DOCKERCOPYFILES="config/config.json" \
			build
			;;
		run)
			DOCKERENV="-e SERVERPORT=8080 -e DESTINATION=http://172.17.0.2:8081" \
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
