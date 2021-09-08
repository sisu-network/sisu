#!/bin/sh

CUR_DIR=$(pwd)
COSMOS_VERSION=0.42.1

for folder in $CUR_DIR/proto/**
do
  dir=$(basename $folder)
  for file in $folder/*
  do
    protoc -I $CUR_DIR/proto/$dir  \
      -I $GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v$COSMOS_VERSION/proto \
      -I $GOPATH/pkg/mod/github.com/cosmos/cosmos-sdk@v$COSMOS_VERSION/third_party/proto \
        --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:./x/$dir/types \
        $file
  done
done

