#/bin/bash

docker pull redis

docker run --name dibbler-redis -p 6379:6379 -d redis 
