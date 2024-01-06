#!/bin/bash
echo 1 > /proc/sys/net/ipv4/ip_forward

service ssh start
service rsyslog start

# Start Go app
/app/main

if [ -z "$@" ]; then
    exec /bin/bash
else
    exec $@
fi
