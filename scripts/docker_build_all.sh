#!/bin/sh

# Build all dockers

cd ../sisu
docker build -f Dockerfile_local -t sisu .

say "Docker build is done!" # In case you are watching Youtube.
