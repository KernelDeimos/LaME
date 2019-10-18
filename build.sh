#!/bin/bash
mkdir -p ./env/go/src/github.com/rosewoodmedia

echo "['log', 'copy-to-env.go.start', 'lamego']"
cp -a ./lamego ./env/go/src/github.com/KernelDeimos/
echo "['log', 'copy-to-env.go.end', 'lamego']"
