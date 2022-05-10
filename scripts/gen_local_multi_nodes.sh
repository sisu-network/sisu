#!/bin/sh

go build -o ./sisu cmd/sisud/main.go
if [ $? -eq 0 ]; then
    echo "Build succeeded"
else
    exit 1
fi

rm -rf output/

./sisu multi-nodes --v 2 --c 2

for ((i = 0; i < 4; i++))
do
  rm -rf ~/.validator"$i"
  cp -rf ./output/validator"$i" ~/.validator"$i"
done

# Copy dheart.toml to its folder
for ((i = 0; i < 4; i++))
do
  mkdir -p ~/.validator"$i"/dheart
  cp ../dheart/dheart.toml ~/.validator"$i"/dheart
done
