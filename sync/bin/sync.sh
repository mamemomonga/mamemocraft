#!/bin/bash
set -eu

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
KEYFILE=../web/etc/id_ed25519
TIME=$(date +'%Y-%m-%d %H:%M:%S')
cd $BASEDIR

echo "$TIME" > log.txt

rsync -av --delete \
	mamemocraft@mc01.mamemo.online:/home/mamemocraft/mamemocraft/data/ data/data/ >> log.txt 2>&1

cd $BASEDIR/data

git add . >> log.txt 2>&1
git commit -a -m "$TIME" >> log.txt 2>&1

