#!/bin/bash

if sudo echo "nameserver 1.1.1.1" >> /etc/resolv.conf
then echo "successfully added nameserver 1.1.1.1 to /etc/resolv.conf"
fi 
exit $?