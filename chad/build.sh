#!/bin/bash
APP=chad
source "$(git rev-parse --show-toplevel)/build_lib.sh"

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
