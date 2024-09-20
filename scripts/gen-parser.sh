#!/bin/bash

echo "Generating parser..."

goyacc -l -o src/gd/ast/gd.y.go src/gd/ast/gd.y

echo "Generation completed!"
