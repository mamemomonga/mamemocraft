#!/bin/bash
set -eux

cp mamemocraft-start.service /etc/systemd/system/
cp mamemocraft-stop.service /etc/systemd/system/

systemctl enable mamemocraft-start.service
systemctl enable mamemocraft-stop.service
systemctl daemon-reload
