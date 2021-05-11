#!/bin/sh

CUR_DIR=$(pwd)
COSMOS_VERSION=0.42.1

for f in $CUR_DIR/proto/evm/*
do
  echo "Processing $f"
protoc -I $CUR_DIR/proto/evm  \
  -I $GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v$COSMOS_VERSION/proto \
  -I $GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v$COSMOS_VERSION/third_party/proto \
    --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:./x/evm/types \
    $f
done