#!/bin/sh

go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc 2>/dev/null

find ./x -name "*.pb.go" -delete

dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $dirs; do
  protoc \
    -I "proto" \
    --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/sisu-network/sisu/* ./
rm -rf github.com
