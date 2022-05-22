#!/bin/bash

go build -o ./sisu cmd/sisud/main.go
if [ $? -eq 0 ]; then
    echo "Build succeeded"
else
    exit 1
fi

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
    --genesis-folder)
      GENESIS_FOLDER="$2"
      shift # past argument
      shift # past value
      ;;
    *)    # unknown option
      POSITIONAL+=("$1") # save it in an array for later
      shift # past argument
      ;;
  esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

genesis_arg=""
if [ -n "$GENESIS_FOLDER" ]
then
  genesis_arg="--genesis-folder $GENESIS_FOLDER"
fi

./sisu localnet $genesis_arg
rm -rf ~/.sisu
cp -rf ./output/node0/ ~/.sisu
rm -rf ./output

# Copy dheart.toml to its folder
mkdir -p ~/.sisu/dheart
cp ../dheart/dheart.toml ~/.sisu/dheart
