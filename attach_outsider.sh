#!/bin/bash

docker rmi kademlia_outsider
docker build -t kademlia_outsider -f ./Dockerfile.outsider .
docker run -it --rm --net kademlia_net1 kademlia_outsider
