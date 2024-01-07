#!/bin/bash

# Set default iptables policies
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# Allow loopback and ICMP
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -p icmp -j ACCEPT

#ssh
INPUT -p tcp --dport 22 -i eth2 -s 10.0.3.3 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 10.0.3.3 -j ACCEPT
iptables -A INPUT -p tcp --sport 22 -s 10.0.3.0/24 -j ACCEPT

#postgres
iptables -A INPUT -p tcp --dport 5432 -j ACCEPT
iptables -A INPUT -p tcp --sport 5432 -j ACCEPT

#https
iptables -A INPUT -p tcp --dport 8080 -s 10.0.1.4 -j ACCEPT
iptables -A INPUT -p tcp --sport 8080 -s 10.0.1.0/24 -j ACCEPT

ip route del default
ip route add default via 10.0.2.2 dev eth0


service ssh start
service rsyslog start

# Start Go app
/app/main

if [ -z "$@" ]; then
    exec /bin/bash
else
    exec $@
fi
