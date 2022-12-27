#!/bin/bash

iptables -A OUTPUT --destination 192.168.110.153 -j NFQUEUE --queue-num 2000
