#!/bin/bash
export LC_NUMERIC="en_US.UTF-8"

APP=nginx
IMAGEBASE=nginx:alpine
NGINXPORT=1080
BASEPATH=$(git rev-parse --show-toplevel)
source "$BASEPATH/build_lib.sh"
BINAB=/cygdrive/c/apache24/bin/ab
DOCKERHOST="192.168.99.100"

usage(){
	printf "$YELLOW* Usage:
	build $BLUE[compile|build|run|clean|list]$RESET
"
}

start_nginx(){
	BASEHASH=$(docker run -d -p $NGINXPORT:80 --name $APP $IMAGEBASE)
}

stop_nginx(){
	docker stop $BASEHASH > /dev/null && docker rm $BASEHASH > /dev/null
}

start_dummy(){
	bash $BASEPATH/dummy/build.sh launch
}

stop_dummy(){
	CONTAINER=$(docker ps --filter "name=dummy" --format "{{.ID}}")
	docker stop $CONTAINER > /dev/null && docker rm $CONTAINER > /dev/null
}

do_test(){
	BASEDATA=$($BINAB -c $1 -n $2 $3 2>&1)
	DATA=$($BINAB -c $1 -n $2 $4 2>&1)
	get_results "$BASEDATA" "$DATA"
}

print_diff(){
	X=$1 && Y=$2
	DIFF=$(bc <<< "scale=9;(($Y/$X)*100)-100")
	POSITIVE=$(echo "$DIFF>=0" | bc -l)
	if [ "$NEGATIVE" = "1" ]
	then
		POSITIVE=$((POSITIVE-2))
	fi
	if [ $POSITIVE -eq 1 -o $POSITIVE -eq -2 ]
	then
		COLOR=$GREEN
	else
		COLOR=$RED
	fi
	printf "   %-21s: $COLOR%+9.4f$RESET   [$1/$2]\n" "$3" $DIFF
}

get_results(){
	BASESPENTTIME=$(echo "$1" | grep "Time taken for tests" | awk '{ print $5 }')
	BASEFAILED=$(echo "$1" | grep "Failed requests" | awk '{ print $3 }')
	BASERPS=$(echo "$1" | grep "Requests per second" | awk '{ print $4 }')
	BASETPR=$(echo "$1" | grep "Time per request:" | head -n 1 | awk '{ print $4 }')
	BASECTPR=$(echo "$1" | grep "Time per request:" | tail -n 1 | awk '{ print $4 }')
	BASETRANSFERRATE=$(echo "$1" | grep "Transfer rate:" | head -n 1 | awk '{ print $3 }')

	SPENTTIME=$(echo "$2" | grep "Time taken for tests" | awk '{ print $5 }')
	FAILED=$(echo "$2" | grep "Failed requests" | awk '{ print $3 }')
	RPS=$(echo "$2" | grep "Requests per second" | awk '{ print $4 }')
	TPR=$(echo "$2" | grep "Time per request:" | head -n 1 | awk '{ print $4 }')
	CTPR=$(echo "$2" | grep "Time per request:" | tail -n 1 | awk '{ print $4 }')
	TRANSFERRATE=$(echo "$2" | grep "Transfer rate:" | head -n 1 | awk '{ print $3 }')

	NEGATIVE=1 print_diff $BASESPENTTIME $SPENTTIME "Spent time"
	# print_diff "$BASEFAILED" "$FAILED" "Failed requests"
	print_diff $BASERPS $RPS "Requests per second"
	NEGATIVE=1 print_diff $BASETPR $TPR "Time per request"
}

URLBASE=${$1:="http://$DOCKERHOST:$NGINXPORT/"}
URL=${$2:="http://$DOCKERHOST:8081/":=}

# printf $YELLOW"+ Starting nginx\n"$RESET
# start_nginx
# printf $YELLOW"+ Starting dummy\n"$RESET
# start_dummy
printf $YELLOW"+ Testing scenario 1\n"$RESET
RESULTS_TEST01=$(do_test 500 2000 $URLBASE $URL)
printf $YELLOW"+ Testing scenario 2\n"$RESET
RESULTS_TEST02=$(do_test 1000 3000 $URLBASE $URL)
# printf $YELLOW"+ Testing scenario 3\n"$RESET
# RESULTS_TEST03=$(do_test 1500 7500 $URLBASE $URL)
# printf -- $YELLOW"- Stoping nginx\n"$RESET
# stop_nginx
# printf -- $YELLOW"- Stoping dummy\n"$RESET
# stop_dummy
printf -- $BLUE"- Scenario 01\n"$RESET
echo "$RESULTS_TEST01"
printf -- $BLUE"- Scenario 02\n"$RESET
echo "$RESULTS_TEST02"
# printf -- $BLUE"- Scenario 03\n"$RESET
# echo "$RESULTS_TEST03"
