#!/bin/bash

# Example usage:
# ./scripts/build-gd-tool.sh go gdvm release darwin amd64 0.0.1

if [ "$#" -ne 8 ]; then
    echo "Usage: $0 [go | tinygo] [gdc | gdcvm | gdvm] [debug | release] [os] [arch] [version] [build_number] [binary_path]"
    exit 1
fi

compiler=$1
tool=$2
build_mode=$3
os=$4
arch=$5
version=$6
build_number=$7
binary_path=$8

ext=""

if [ "$os" == "windows" ]; then
    ext=".exe"
fi

if [ "$os" == "js" ]; then
    ext=".wasm"
fi

binary_file_name=$tool$ext

echo "Building [$build_mode] | Tool: $tool | OS: $os | Arch: $arch | Version: $version | Build: $build_number"

export GOOS=$os
export GOARCH=$arch
export CGO_ENABLED=0

project_path="./src/cmd/$tool/main.go"

echo "Compiling: $project_path"

if [ "$build_mode" == "debug" ]; then
    # Development build (with debugging information)
    $compiler build -tags debug -o $binary_path/$binary_file_name $project_path
elif [ "$build_mode" == "release" ]; then
    ldflags="-X 'main.version=$version' -X 'main.buildNumber=$build_number' -X 'main.arch=$arch'"
    # Release build (optimized and without debugging information)
    if [ "$compiler" == "go" ]; then
        go build -o $binary_path/$binary_file_name -trimpath -buildvcs=false -ldflags "-s -w $ldflags" $project_path
    elif [ "$compiler" == "tinygo" ]; then
        tinygo build -o $binary_path/$binary_file_name -opt=z -panic=trap -no-debug -scheduler=none -ldflags "$ldflags" $project_path
    fi
else
    echo "Error: Unrecognized build mode: $build_mode"
    exit 1
fi

echo "Build successful! Output: $binary_path/$binary_file_name"
