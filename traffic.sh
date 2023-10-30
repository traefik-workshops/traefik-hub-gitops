#!/usr/bin/env bash
#
TOKEN=PQRbrvwK46kY39gMJ9jpwf1iTgaHvxfeKZJA5QdR
GW_DOMAIN="subjective-squirrel-ktmkru.cv5yg484.traefikhub.io"

random() {
	NUMBER=$RANDOM
	let "NUMBER%=$1"

	# add 10 to have a minimum of 10 for the -c option of hey
	let "NUMBER+=10"
	echo $NUMBER
}

call() {
   METHOD=$1
   HPATH=$2
   RDNB=$3
   ERRORRATE=$4
   NB=$(random $RDNB)

   echo "$HPATH => "

   let "ERR=$NB/100*$ERRORRATE"

   hey -m "$METHOD" -c 10 -n $NB \
  	-H 'accept: application/json' \
  	-H "Authorization: Bearer $TOKEN" "https://$GW_DOMAIN$HPATH" | grep responses
   hey -m "$METHOD" -c 1 -n $ERR \
      	-H 'accept: application/json' \
      	-H "Authorization: Bearer $TOKEN" "https://$GW_DOMAIN$HPATH?status=503" | grep responses
       echo ""
   echo ""

}

for (( ; ; ))
do
 call 'GET' "/customers/customers" 1000 2
 call 'GET' "/flights/flights" 1000 5
 call 'GET' "/tickets/tickets" 1000 2
 call 'POST' "/tickets/tickets" 1000 2
 call 'POST' "/employees/employees" 1000 2

 sleep 2;
done;
