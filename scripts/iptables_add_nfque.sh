#!/bin/bash

# docker exec -it wade-test-1 bash,then excute follow cmd
#iptables -A OUTPUT --destination 192.168.110.153 -j NFQUEUE --queue-num 2000
#iptables -t raw -I PREROUTING -p tcp --syn -j NFQUEUE --queue-num=1 --queue-bypass
iptables -t raw -I OUTPUT -p tcp --syn -j NFQUEUE --queue-num=1 --queue-bypass
