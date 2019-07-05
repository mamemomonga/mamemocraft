#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

rsync -av --delete \
	-e "ssh -i web/etc/id_ed25519" \
	--exclude='spigot/*' \
	--exclude='data/*' \
	--exclude='var/*' \
	--exclude='README.md' \
	mamemocraft@mc01.mamemo.online:/home/mamemocraft/mamemocraft/ \
	mamemocraft/

