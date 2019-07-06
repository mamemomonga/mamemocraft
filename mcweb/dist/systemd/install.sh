#!/bin/bash
set -eux
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../.. && pwd )"

cat mamemocraft-web.service | perl -nlp -E 's!###BASEDIR###!'$BASEDIR'!g' > /etc/systemd/system/mamemocraft-web.service

systemctl daemon-reload
systemctl enable mamemocraft-web
systemctl start mamemocraft-web
systemctl status mamemocraft-web

