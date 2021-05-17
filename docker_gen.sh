#!/bin/sh

# Generate binary file
rm -rf build
mkdir -p build/bin
cd src
go build -o=./../build/bin/sisu ./cmd/sisud/main.go
cp .env.dev ../build/.env
cd ..

# Generate genesis data for testnet
OUTPUT=output
NODE_COUNT=2

cd build
# Generate nodes
./bin/sisu testnet --v $NODE_COUNT --starting-ip-address 192.168.10.2

for (( i=0; i<$NODE_COUNT; i++ ))
do
  # Disable strict address book to allow private IPs
  sed -i '' -e 's/addr_book_strict = true/addr_book_strict = false/' ./$OUTPUT/node$i/main/config/config.toml
done

cd ..