#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. &&  pwd )"
source $BASEDIR/config

exec $BASEDIR/bin/mcrcon -H localhost -p $RCON_PASSWORD "say サーバを停止します" save-all stop
