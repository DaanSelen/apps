#!/bin/bash

touch ./config.txt
echo -e '# NMAC CONFIG FILE\n
# Nerthus Monitoring Application Server IP Address.
manager_ip=\n
# Set to true or false if you want uptime enabled or not.
uptime_on=true\n
# (Optional) If desired, place custom lockpath directory here. (WIHTOUT TRAILING /) Make sure the program has access to the location.
lock_dir=\n
# (Optional) If Desired, set the monitoring interval in seconds.
mon_interval=' > ./config.txt