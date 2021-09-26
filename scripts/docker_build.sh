#!/bin/sh

cp ~/.ssh/id_rsa ./tmp
docker build . -t sisu
docker run --name=sisu -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -t sisu

# Clean up the container
rm ./build/id_rsa
# docker rm -f sisu
# docker rmi -f sisu