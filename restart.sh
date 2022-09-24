#!/bin/bash

docker compose down && \
docker rmi docker-go | true && \
docker build -t docker-go . && \
docker compose up -d
