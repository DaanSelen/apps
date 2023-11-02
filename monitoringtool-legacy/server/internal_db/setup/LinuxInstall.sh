#!/bin/bash

if [[ ! $(sudo echo 0) ]]; then echo "Please run as sudo"; fi
chmod 777 ../CSMTS
cp ../CSMTS.service /lib/systemd/system
systemctl enable CSMTS && systemctl restart CSMTS
echo "Script Finished, Showing Status:"
systemctl status CSMTS