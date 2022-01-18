#!/bin/sh

# Build all dockers

cd ../deyes
docker build -t deyes .
docker build -f Dockerfile-ganache -t ganache .

cd ../dheart
docker build -t dheart .

cd ../sisu
docker build -f Dockerfile_local -t sisu .

say "Docker build is done!" # In case you are watching Youtube.
