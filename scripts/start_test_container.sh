#!/bin/bash

docker run -d --privileged=true --name wade-test-1 wade23/deploy:deploytest
docker run -d --privileged=true --name wade-test-2 wade23/deploy:deploytest

