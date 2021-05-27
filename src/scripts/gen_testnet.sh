#!/bin/sh

go build -o ./sisu cmd/sisud/main.go
./sisu localnet
rm -rf ~/.sisu
cp -rf ./output/node0/ ~/.sisu
rm -rf ./output