#!/bin/sh

# Debug: docker run -it sisu ls /dist

if (( $# < 1 )); then
  echo 'Please provide the node name'
  exit 1
fi

echo Running $1

docker run --rm -v $(pwd)/build/output/$1:/root/.sisu sisu /dist/sisu start