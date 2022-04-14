#!/bin/sh

export PRIVATE_KEY=$(cat actions-observer-local.2022-04-11.private-key.pem)
docker run -it --rm --name github-prometheus-client -p 8080:8080 -e PRIVATE_KEY --env-file=.env.docker github-prometheus-client	