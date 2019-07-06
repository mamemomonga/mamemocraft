#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

rsync -av --delete \
	-e "ssh -i mcweb/etc/id_ed25519" \
	--exclude='spigot/*' \
	--exclude='.gitignore' \
	--exclude='data/*' \
	--exclude='data.bak/*' \
	--exclude='var/*' \
	--exclude='README.md' \
	mamemocraft@mc01.mamemo.online:/home/mamemocraft/mamemocraft/ \
	mamemocraft/

