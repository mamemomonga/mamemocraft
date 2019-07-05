#!/bin/bash
set -eux

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

exec autossh -i controller/etc/id_ed25519 -L 5006:localhost:8123 -N mamemocraft@mc01.mamemo.online

