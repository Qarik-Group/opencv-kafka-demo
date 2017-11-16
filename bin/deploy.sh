#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd $ROOT

cf push "$@"
