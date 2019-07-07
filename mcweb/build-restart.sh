#!/bin/bash
set -eux
make
sudo systemctl restart mamemocraft-web
exec sudo journalctl -u mamemocraft-web -f

