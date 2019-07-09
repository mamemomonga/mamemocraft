#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

make
sudo systemctl restart mamemocraft-web
exec sudo journalctl -u mamemocraft-web -f

