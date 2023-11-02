#!/bin/bash

if [[ ! $(sudo echo 0) ]]; then echo "Please run as sudo"; fi
chmod 760 ../CSMTA
cp ../CSMTA.service /lib/systemd/system
systemctl enable CSMTA && systemctl restart CSMTA
echo "Script Finished, Showing Status:"
systemctl status CSMTA