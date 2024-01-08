#!/bin/bash

touch ./config.txt
echo -e '
# NMAC CONFIG FILE\n\n
# Nerthus Monitoring Application Server IP Address.\n
manager_ip=\n\n
# (Optional) If Desired, Place Custom Lockpath Directory Here. (WIHTOUT TRAILING /)\n
lock_dir=' > ./config.txt