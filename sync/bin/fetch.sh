#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
KEYFILE=../web/etc/id_ed25519
cd $BASEDIR

rsync -av --delete \
	mamemocraft@mc01.mamemo.online:/home/mamemocraft/mamemocraft/data/ data/data/

