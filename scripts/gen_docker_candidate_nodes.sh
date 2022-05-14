#!/bin/sh

go build -o ./sisu cmd/sisud/main.go
if [ $? -eq 0 ]; then
    echo "Build succeeded"
else
    exit 1
fi

rm -rf output/

./sisu local-docker --v 2 --c 1 --output-dir ./output
