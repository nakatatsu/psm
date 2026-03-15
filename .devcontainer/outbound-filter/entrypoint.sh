#!/bin/sh
chmod 666 /proc/1/fd/1 /proc/1/fd/2
exec squid -N -d 1
