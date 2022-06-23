#!/bin/bash

set -e

cd "$(dirname "$0")"

rm -rf ./output
mkdir ./output

# go build
for dir in ./cmd/*; do
    pkg=$( basename "$dir" )
    go build -v -o ./output/"$pkg" ./cmd/"$pkg"
done