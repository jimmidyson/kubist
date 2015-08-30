#!/bin/bash

FILES=$(gofmt -l `find . -type f -name "*.go" | grep -v Godeps`)
if [[ ! -z "$FILES"  ]]; then
  echo Run gofmt on the following files:$'\n' $FILES
  exit 1
fi
