#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
MINECRAFT_SERVER=1.14.3
SPIGOT_DIR=$BASEDIR/spigot

mkdir -p $SPIGOT_DIR
cd $SPIGOT_DIR
curl -L -o BuildTools.jar https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar
java -jar BuildTools.jar --rev $MINECRAFT_SERVER
