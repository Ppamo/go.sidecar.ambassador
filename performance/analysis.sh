#!/bin/bash
export LC_NUMERIC="en_US.UTF-8"

APP=nginx
IMAGEBASE=nginx:1.13.9-alpine
NGINXPORT=1080
BASEPATH=$(git rev-parse --show-toplevel)
source "$BASEPATH/build_lib.sh"

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
	go get github.com/rakyll/hey > /dev/null 2>&1
	BASEDATA=$(hey -c $1 -n $2 $3)
	DATA=$(hey -c $1 -n $2 $4)
	get_results "$BASEDATA" "$DATA"
}

print_diff(){
	DIFF=$(bc <<< "scale=9;(($1/$2)*100)-100")
	POSITIVE=$(echo "$DIFF>=0" | bc -l)
	if [ "$NEGATIVE" = "1" ]
	then
		POSITIVE=$((POSITIVE-1))
	fi
	if [ $POSITIVE -eq 1 ]
	then
		COLOR=$GREEN
	else
		COLOR=$RED
	fi
	printf "\t%-21s: $COLOR%+9.4f \n$RESET" "$3" $DIFF
}

get_results(){
	BASETOTAL=$(echo "$1" | grep "Total:" | awk '{ print $2 }')
	BASESLOWEST=$(echo "$1" | grep "Slowest:" | awk '{ print $2 }')
	BASEFASTEST=$(echo "$1" | grep "Fastest:" | awk '{ print $2 }')
	BASEAVERAGE=$(echo "$1" | grep "Average:" | awk '{ print $2 }')
	BASETPS=$(echo "$1" | grep "Requests/sec:" | awk '{ print $2 }')

	TOTAL=$(echo "$2" | grep "Total:" | awk '{ print $2 }')
	SLOWEST=$(echo "$2" | grep "Slowest:" | awk '{ print $2 }')
	FASTEST=$(echo "$2" | grep "Fastest:" | awk '{ print $2 }')
	AVERAGE=$(echo "$2" | grep "Average:" | awk '{ print $2 }')
	TPS=$(echo "$2" | grep "Requests/sec:" | awk '{ print $2 }')

	printf "\nResults:\n"
	NEGATIVE=1 print_diff $BASETOTAL $TOTAL "Spent time"
	NEGATIVE=1 print_diff $BASESLOWEST $SLOWEST "Slowest response"
	NEGATIVE=1 print_diff $BASEFASTEST $FASTEST "Fastest response"
	NEGATIVE=1 print_diff $BASEAVERAGE $AVERAGE "Average response"
	print_diff $BASETPS $TPS "Responses per second"
}

printf $YELLOW"+ Starting nginx\n"$RESET
start_nginx
printf $YELLOW"+ Starting dummy\n"$RESET
start_dummy
printf $YELLOW"+ Testing scenario 1\n"$RESET
RESULTS=$(do_test 100 500 http://localhost:$NGINXPORT http://localhost:8081)
printf -- $YELLOW"- Stoping nginx\n"$RESET
stop_nginx
printf -- $YELLOW"- Stoping dummy\n"$RESET
stop_dummy
echo "$RESULTS"
