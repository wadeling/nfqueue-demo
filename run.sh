#!/bin/bash
setcap 'cap_net_admin=+ep' 
./nfqueue-demo
