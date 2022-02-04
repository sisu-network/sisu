#!/bin/sh

SISU_COUNT="$1"

# Build the local sisu
cd ..
go build -o ./sisu cmd/sisud/main.go

# generate the testnet data
./sisu testnet --v $SISU_COUNT --output-dir ./output --config-string "$(
  echo '{'
  echo '"chains":[{"name":"ganache1","rpc":"http://ganache1:7545"},{"name":"ganache2","rpc":"http://ganache2:7545"}],'
  echo '"nodes":['
  for i in $(seq 1 $SISU_COUNT); do
    if [ "$i" != 1 ]; then echo ','; fi
    printf '{"sisu_ip": "sisud.sisu-%d", "dheart_ip":"dheart.sisu-%d","deyes_ip":"deyes.sisu-%d","sql":{"host":"mysql.mysql","port":3306,"username":"root","password":"password"}}' $i $i $i
  done
  echo ']'
  echo '}'
)"

# move final output to kind
rm -rf ./kind/output
mv output kind
