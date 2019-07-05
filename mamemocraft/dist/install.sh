#!/bin/bash
set -eux

cp mamemocraft-start.service /etc/systemd/system/
cp mamemocraft-stop.service /etc/systemd/system/

systemctl enable mamemocraft-start.service
systemctl enable mamemocraft-stop.service
systemctl daemon-reload


cp sudoers /etc/sudoers.d/mamemocraft
chown root:root /etc/sudoers.d/mamemocraft
chmod 600 /etc/sudoers.d/mamemocraft
