#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
cd $BASEDIR

ssh -i controller/etc/id_ed25519 mamemocraft@mc01.mamemo.online tar zcC /home/mamemocraft mamemocraft/bin mamemocraft/systemd | tar zxv

