#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. &&  pwd )"

exec $BASEDIR/bin/mcrcon -H localhost -p minecraft "say サーバを停止します" save-all stop
