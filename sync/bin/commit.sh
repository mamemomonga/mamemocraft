#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
KEYFILE=../web/etc/id_ed25519
cd $BASEDIR/data

TIME=$(date +'%Y-%m-%d %H:%M:%S')

git add .
git commit -a -m "$TIME"

