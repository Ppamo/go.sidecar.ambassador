#!/bin/bash
APP=mockserver
OS=$(uname -o)
source "$(git rev-parse --show-toplevel)/lib/build_lib.sh"

BASEPATH="$(git rev-parse --show-toplevel)/mockserver"
BASEPATH="/c${BASEPATH#C:}"

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

CLEANBASEPATH=$BASEPATH
if [ "$OS" == "Cygwin" ]
then
	CLEANBASEPATH=${BASEPATH#/cygdrive}
fi
while [ $INDEX -lt $# ]
do
	((INDEX++))
	case ${@:INDEX:1} in
		compile)	compile			;;
		build)		build			;;
		run)
			DOCKERVOLUMES="-v $CLEANBASEPATH/responses:/mocks" \
			run
			;;
		launch)		launch			;;
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
