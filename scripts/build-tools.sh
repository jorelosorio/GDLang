#!/bin/bash

# Example usage:
# ./scripts/build-tools.sh go release darwin amd64 0.0.1

if [ "$#" -ne 8 ]; then
    echo "Usage: $0 [go | tinygo] [debug | release] [os] [arch] [version] [build_number] [binary_path] [dist_path]"
    exit 1
fi

# Get the parameter value
compiler=$1
build_mode=$2
os=$3
arch=$4
version=$5
build_number=$6
binary_path=$7
dist_path=$8

rm -rf $binary_path
mkdir -p $binary_path

# List of tools
tools=("gdc" "gdcvm" "gdvm")

for tool in "${tools[@]}"; do
    ./scripts/build-gd-tool.sh "$compiler" "$tool" "$build_mode" "$os" "$arch" "$version" "$build_number" $binary_path
done

mkdir -p "$dist_path"

echo "🚀 Compressing..."

# Compress the tool binaries into a tarball
tar --exclude='./.DS_Store' \
    --exclude='./*.DS_Store' \
    --exclude='./.*' \
    --exclude='./*.swp' \
    --exclude='./*.swo' \
    --exclude='./.git' \
    -czvf "$dist_path/gdlang-$version-$os-$arch.tar.gz" -C $binary_path ./ -C ../ LICENSE.txt

echo "🎉 Done!"
