#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

rm -f var/running
while true; do
	if [ -z "$(grep '\[Server thread\/INFO\]: Done' var/log.txt)" ]; then
		sleep 1
	else
		touch var/running
		break
	fi
done

