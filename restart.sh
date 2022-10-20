#!/bin/bash

docker compose down && \
docker rmi docker-go | true && \
docker build -t docker-go . && \
docker compose up -d

if [ "$1" = "-a" ]
  then
    ./attach_outsider.sh
else
    echo 'Did not attach any outsider node. Use "./restart.sh -a" to attach one.'
fi
