#!/bin/bash

FILES=$(go vet github.com/fabric8io/kubist/...)
if [[ ! -z "$FILES"  ]]; then
  echo Run go fix on the following files:$'\n' $FILES
  exit 1
fi
