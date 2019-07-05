#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. &&  pwd )"
MINECRAFT_SERVER=1.14.3
MEMSIZE=10G

$BASEDIR/bin/launch_watcher.sh &
cd $BASEDIR/data

rm -f $BASEDIR/var/down

set +e
java -Xms$MEMSIZE -Xmx$MEMSIZE \
	-XX:+AlwaysPreTouch \
	-XX:+DisableExplicitGC \
	-XX:+UseG1GC \
	-XX:+UnlockExperimentalVMOptions \
	-XX:MaxGCPauseMillis=50 \
	-XX:G1HeapRegionSize=4M \
	-XX:TargetSurvivorRatio=90 \
	-XX:G1NewSizePercent=50 \
	-XX:G1MaxNewSizePercent=80 \
	-XX:InitiatingHeapOccupancyPercent=10 \
	-XX:G1MixedGCLiveThresholdPercent=50 \
 	-jar ../spigot/spigot-$MINECRAFT_SERVER.jar \
	> $BASEDIR/var/log.txt 2>&1

touch $BASEDIR/var/down

