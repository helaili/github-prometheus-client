#!/bin/sh
#export PRIVATE_KEY=$(cat actions-observer-local.2022-04-11.private-key.pem)
docker run -d --name redis-stack -e REDIS_ARGS="--requirepass password" -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
docker run -it --rm --name github-prometheus-client -p 8080:8080 --env-file=.env.docker github-prometheus-client	
